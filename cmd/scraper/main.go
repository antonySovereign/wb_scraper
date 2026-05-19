package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"wb_scraper/internal/config"
	"wb_scraper/internal/repository"
	"wb_scraper/internal/scraper"
	"wb_scraper/internal/service"
)

const (
	baseURL = "https://wildberries.ru"
	workers = 2
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// connect to postgres
	db, err := repository.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get db connection: %w\n", err)
	}

	// connect to redis
	redisClient, err := repository.InitRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Successfully connected to redis...")

	// connect to kafka
	if err := repository.CheckKafkaConnection([]string{cfg.KafkaBrokers}); err != nil {
		log.Fatal(err)
	}

	kafkaWriter := repository.InitKafkaWriter([]string{cfg.KafkaBrokers}, cfg.KafkaTopic)
	log.Print("Successfully connected to kafka...")
	defer kafkaWriter.Close()

	categoryRepo := repository.NewCategoryRepo(db)
	categoryService := service.NewCategoryService(categoryRepo)

	productsRepo := repository.NewProductsRepository(db)
	productsService := service.NewProductService(productsRepo)

	// root context for all operations, can be used to cancel all operations if needed
	rootCtx, cancelAll := context.WithCancel(ctx)
	defer cancelAll()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, exiting...")
		cancelAll()
	}()

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

			for {
				select {
				case <-rootCtx.Done():
					log.Printf("Worker %d: shutting down\n", workerID)
					return
				case job, ok := <-jobs:
					if !ok {
						log.Printf("Worker %d: no more jobs, exiting\n", workerID)
						return
					}

					log.Printf("Worker %d: parsing %s", workerID, job.URL)

					productsJson, err := localScr.GetProducts(job.URL)
					if err != nil {
						log.Printf("Worker got error %d: %v", workerID, err)
						continue
					}

					if err := productsService.SyncProducts(ctx, productsJson, job.CategoryID); err != nil {
						log.Printf("Worker got error %d: %v", workerID, err)
						continue
					}
				}
			}
		}(i)
	}

	for _, category := range categories {
		select {
		case <-rootCtx.Done():
			log.Println("Main: shutting down, no more jobs will be added")
			break
		case jobs <- Job{CategoryID: category.DbID, URL: baseURL + category.Url}:
		}
	}

	close(jobs)

	wg.Wait()
	log.Print("All workers finished!!!")

	log.Print("Closing db connection...")
	if err := sqlDB.Close(); err != nil {
		log.Printf("Failed to close db connection: %v", err)
	} else {
		log.Print("Db connection closed successfully")
	}

	log.Print("Closing redis connection")
	if err := redisClient.Close(); err != nil {
		log.Printf("Failed to close redis connection: %v", err)
	} else {
		log.Print("Redis connection closed successfully")
	}

	log.Print("Scraper finished")
}
