package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"wb_scraper/internal/config"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: cfg.DBSchema + ".",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to db: %w", err)
	}

	return db, nil

}
