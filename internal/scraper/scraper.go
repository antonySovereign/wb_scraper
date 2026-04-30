package scraper

import (
	"context"
	"fmt"
	"strings"
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
