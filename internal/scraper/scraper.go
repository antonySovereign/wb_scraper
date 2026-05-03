package scraper

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

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

func (s *Scraper) GetProducts(url string) error {
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

	dataChan := make(chan string, 10)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if responseEv, ok := ev.(*network.EventResponseReceived); ok {
			if strings.Contains(responseEv.Response.URL, "https://www.wildberries.ru/__internal/u-search/") {
				// Выполняем получение тела в отдельном контексте, чтобы не вешать основной цикл
				go func(id network.RequestID) {
					var body []byte
					// Важно: используем chromedp.Run с базовым контекстом
					err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
						var err error
						body, err = network.GetResponseBody(id).Do(ctx)
						return err
					}))

					if err == nil {
						dataChan <- string(body)

						path := responseEv.Response.URL + ".json"
						file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

						if err != nil {
							log.Print(err)
						} else {
							defer file.Close()
							writer := bufio.NewWriter(file)
							_, err = writer.Write(body)
							if err == nil {
								err = writer.Flush()
								if err == nil {
									log.Printf("Written to file: %s", path)
								}
							}

						}
					} else {
						log.Print(err)
					}

				}(responseEv.RequestID)
			}
		}
	})

	return chromedp.Run(ctx,
		network.Enable(),
		// Эмулируем поведение человека: сначала на главную (чтобы получить куки)
		chromedp.Navigate("https://www.wildberries.ru/"),
		chromedp.Sleep(7*time.Second),

		// Переходим на страницу категории
		chromedp.Navigate(url),

		// Ожидаем появления контента, чтобы триггернуть запрос к API
		chromedp.WaitVisible(`.catalog-page`, chromedp.ByQuery),

		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				select {
				case jsonBody := <-dataChan:
					if strings.Contains(jsonBody, `"products"`) {
						fmt.Printf("Успех! Данные получены. Размер: %d символов\n", len(jsonBody))
						return nil
					}
				case <-ctx.Done():
					return fmt.Errorf("превышено время ожидания ответа от API")
				}
			}
		}),
	)
}
