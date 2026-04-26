package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	rawJSON := []byte(`[{"id":131289,"name":"Экспресс","url":"/catalog/ekspress-dostavka"},{"id":306,"name":"Женщинам","url":"/catalog/zhenshchinam","childs":[{"id":8126,"name":"Блузки и рубашки"}]}]`)

	var category []Category
	if err := json.Unmarshal(rawJSON, &category); err != nil {
		log.Fatalf("Failed to parse json: %v", err)
	}

	fmt.Println(len(category))

	// opts := append(chromedp.DefaultExecAllocatorOptions[:],
	// 	chromedp.NoFirstRun,
	// 	chromedp.NoDefaultBrowserCheck,
	// 	chromedp.Flag("headless", false),
	// 	chromedp.Flag("disable-blink-features", "AutomationControlled"),
	// )

	// alloCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	// defer cancel()

	// ctx, cancel := chromedp.NewContext(alloCtx)
	// defer cancel()

	// menuURLChan := make(chan string, 1)

	// chromedp.ListenTarget(ctx, func(ev interface{}) {
	// 	if ev, ok := ev.(*network.EventRequestWillBeSent); ok {
	// 		url := ev.Request.URL
	// 		if strings.Contains(url, "main-menu-ru-ru") {
	// 			menuURLChan <- url
	// 		}
	// 	}
	// })

	// err := chromedp.Run(ctx,
	// 	network.Enable(),
	// 	network.SetCacheDisabled(true),
	// 	chromedp.ActionFunc(func(ctx context.Context) error {
	// 		err := runtime.AddBinding("navigator.driver").Do(ctx)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		_, _, err = runtime.Evaluate(`Object.defineProperty(navigator, 'webdriver', {get: () => undefined})`).Do(ctx)
	// 		return err
	// 	}),
	// 	chromedp.Navigate("https://www.wildberries.ru/"),
	// 	chromedp.WaitVisible(`.nav-element__burger`, chromedp.ByQuery),
	// )

	// if err != nil {
	// 	log.Fatal("Something went wrong..", "err", err)
	// }

	// select {
	// case url := <-menuURLChan:
	// 	fmt.Printf("Url fetched: %s\n", url)
	// case <-time.After(time.Second * 10):
	// 	fmt.Println("Failed to fetch the link")
	// }
}
