package scraper

import (
	"context"
	"fmt"

	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type Scraper struct {
	Headless             bool
	DisableBlinkFeatures string
	UserAgent            string
}

func NewScraper(headless bool, disableBlinkFeatures string, userAgent string) *Scraper {
	return &Scraper{
		Headless:             headless,
		DisableBlinkFeatures: disableBlinkFeatures,
		UserAgent:            userAgent,
	}
}

func (s *Scraper) GetMenuURL() (string, error) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", s.Headless),
		chromedp.Flag("disable-blink-features", s.DisableBlinkFeatures),
		chromedp.UserAgent(s.UserAgent),
	)

	alloCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(alloCtx)
	defer cancel()

	menuURLChan := make(chan string, 10)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*network.EventRequestWillBeSent); ok {
			url := ev.Request.URL
			if strings.Contains(url, "main-menu-ru-ru") {
				menuURLChan <- url
			}
		}
	})

	err := chromedp.Run(ctx,
		network.Enable(),
		network.SetCacheDisabled(true),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := runtime.AddBinding("navigator.driver").Do(ctx)
			if err != nil {
				return err
			}

			_, _, err = runtime.Evaluate(`Object.defineProperty(navigator, 'webdriver', {get: () => undefined})`).Do(ctx)
			return err
		}),
		chromedp.Navigate("https://www.wildberries.ru/"),
		chromedp.WaitVisible(`.nav-element__burger`, chromedp.ByQuery),
	)

	if err != nil {
		return "", err
	}

	select {
	case url := <-menuURLChan:
		return url, nil
	case <-time.After(time.Second * 10):
		return "", fmt.Errorf("Failed to fetch the link")
	}
}

func (s *Scraper) GetProducts(url string) (string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", s.Headless),
		chromedp.Flag("disable-blink-features", s.DisableBlinkFeatures),
		chromedp.UserAgent(s.UserAgent),
	)

	alloCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := context.WithTimeout(alloCtx, 60*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	dataChan := make(chan string, 20)
	var mu sync.Mutex
	pendingRequests := make(map[network.RequestID]bool)
	var result string

	chromedp.ListenTarget(ctx, func(ev interface{}) {

		switch e := ev.(type) {
		case *network.EventResponseReceived:
			if strings.Contains(e.Response.URL, "__internal/u-search/") {
				mu.Lock()
				pendingRequests[e.RequestID] = true
				mu.Unlock()
			}

		case *network.EventLoadingFinished:

			mu.Lock()
			isPending := pendingRequests[e.RequestID]
			if isPending {
				delete(pendingRequests, e.RequestID)
			}
			mu.Unlock()

			if isPending {
				go func(id network.RequestID) {
					var body []byte
					err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
						var err error
						body, err = network.GetResponseBody(id).Do(ctx)
						return err
					}))

					if err == nil {

						select {
						case dataChan <- string(body):
						case <-time.After(time.Second):
						}
					}
				}(e.RequestID)
			}
		}
	})

	err := chromedp.Run(ctx,
		network.Enable(),

		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, err := runtime.Evaluate(`Object.defineProperty(navigator, 'webdriver', {get: () => undefined})`).Do(ctx)
			return err
		}),
		chromedp.Navigate(url),

		chromedp.WaitVisible(`.catalog-page`, chromedp.ByQuery),

		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				select {
				case jsonBody := <-dataChan:
					if strings.Contains(jsonBody, `"products"`) {
						result = jsonBody
						return nil
					}
				case <-ctx.Done():

					return ctx.Err()
				}
			}
		}),
	)

	if err != nil {
		return "", err
	}

	return result, nil
}
