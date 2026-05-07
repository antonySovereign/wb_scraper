package main

import (
	"context"
	"log"

	"wb_scraper/internal/config"
	"wb_scraper/internal/repository"
	"wb_scraper/internal/scraper"
	"wb_scraper/internal/service"
)

func main() {

	// todo LOAD CONFIG

	cfg := config.Load()

	// todo Connect to db

	db, err := repository.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	categoryRepo := repository.NewCategoryRepo(db)
	categoryService := service.NewCategoryService(categoryRepo)

	// productsRepo := repository.NewProductsRepository(db)
	// productsService := service.NewProductService(productsRepo)

	scr := scraper.NewScraper(
		cfg.ChromedpHeadless,
		cfg.ChromedpDisableBlinkFeatures,
		cfg.ChromedpUserAgent,
	)

	url, err := scr.GetMenuURL()
	if err != nil {
		log.Fatal("Failed to parse menu: %w\n", err)
	}

	if err := categoryService.SyncCategories(context.Background(), url); err != nil {
		log.Fatal("Failed to parse and write to db: %w", err)
	}

	// productsJson, err := scr.GetProducts("https://www.wildberries.ru/catalog/muzhchinam/odezhda/dzhinsy")
	// log.Print(productsJson, err)

}
