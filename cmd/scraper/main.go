package main

import (
	"log"

	"wb_scraper/internal/config"
	"wb_scraper/internal/scraper"
)

func main() {

	// todo LOAD CONFIG

	cfg := config.Load()

	// todo Connect to db

	// db, err := repository.InitDB(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// repo := repository.NewCategoryRepo(db)
	// svc := service.NewCategoryService(repo)

	// todo parse menu (test)
	scr := scraper.NewScraper(
		cfg.ChromedpHeadless,
		cfg.ChromedpDisableBlinkFeatures,
		cfg.ChromedpUserAgent,
	)

	// url, err := scr.GetMenuURL()
	// if err != nil {
	// 	log.Fatal("Failed to parse menu: %w\n", err)
	// }

	// if err := svc.SyncCategories(context.Background(), url); err != nil {
	// 	log.Fatal("Failed to parse and write to db: %w", err)
	// }

	log.Fatal(scr.GetProducts("https://www.wildberries.ru/catalog/muzhchinam/odezhda/dzhinsy"))

}
