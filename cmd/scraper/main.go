package main

import (
	"context"
	"log"
	"sync"

	"wb_scraper/internal/config"
	"wb_scraper/internal/repository"
	"wb_scraper/internal/scraper"
	"wb_scraper/internal/service"
)

const (
	baseURL = "https://wildberries.ru"
	workers = 10
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := repository.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	categoryRepo := repository.NewCategoryRepo(db)
	categoryService := service.NewCategoryService(categoryRepo)

	productsRepo := repository.NewProductsRepository(db)
	productsService := service.NewProductService(productsRepo)

	scr := scraper.NewScraper(
		cfg.ChromedpHeadless,
		cfg.ChromedpDisableBlinkFeatures,
		cfg.ChromedpUserAgent,
	)

	// 1. Update categories
	menuURL, err := scr.GetMenuURL()
	if err != nil {
		log.Fatal("Failed to parse menu: %w\n", err)
	}

	if err := categoryService.SyncCategories(ctx, menuURL); err != nil {
		log.Fatal("Failed to parse and write to db: %w", err)
	}

	// 2. Get categories
	categories, err := categoryRepo.GetAll(ctx)
	if err != nil {
		log.Fatal("Failed to get categories from bd: %w", err)
	}

	log.Printf("Found categories: %d", len(categories))

	type Job struct {
		CategoryID int
		URL        string
	}
	jobs := make(chan Job, len(categories))

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			localScr := scraper.NewScraper(
				cfg.ChromedpHeadless,
				cfg.ChromedpDisableBlinkFeatures,
				cfg.ChromedpUserAgent,
			)

			for job := range jobs {
				log.Printf("Worker %d: parsing %s", workerID, job.URL)

				productsJson, err := localScr.GetProducts(job.URL)
				if err != nil {
					log.Printf("Worker got error %d: %v", workerID, err)
					continue
				}

				productsService.SyncProducts(ctx, productsJson, job.CategoryID)
			}
		}(i)
	}

	for _, category := range categories {
		jobs <- Job{CategoryID: category.DbID, URL: baseURL + category.Url}
	}

	close(jobs)

	wg.Wait()
	log.Print("[+] WB Parser finished!!!")

	// for index, category := range categories {
	// 	log.Printf("[%d] Working on category: %s\n", index, category.Name)

	// 	fullURL := baseURL + category.Url

	// 	productsJSON, err := scr.GetProducts(fullURL)
	// 	if err != nil {
	// 		log.Printf("! Failed to parse %s: %v", category.Name, err)
	// 		continue
	// 	}

	// 	if err := productsService.SyncProducts(ctx, productsJSON, category.DbID); err != nil {
	// 		log.Printf("! Failed to sync %s: %v", category.Name, err)
	// 		continue
	// 	}

	// 	log.Printf("[*] Successfuly updated: %s", category.Name)
	// 	time.Sleep(2 * time.Second)

	// }
}
