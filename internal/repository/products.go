package repository

import (
	"context"
	"log"
	"wb_scraper/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (repo *ProductsRepository) Create(ctx context.Context, product *domain.Product) error {
	return repo.db.Table("wb_scraper.products").WithContext(ctx).Create(product).Error
}

func (repo *ProductsRepository) Delete(ctx context.Context, id int) error {
	res := repo.db.Table("wb_scraper.products").
		WithContext(ctx).
		Where("wb_id = ?", id).
		Delete(&domain.Product{})
	if res.Error == nil && res.RowsAffected == 0 {
		log.Printf("Warning: product with id [%d] not found (nothing to delete)", id)
	}
	return res.Error
}

func (repo *ProductsRepository) SaveBatch(ctx context.Context, products []domain.Product) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, product := range products {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "wb_id"}},
				UpdateAll: true,
			}).Omit("Sizes").Create(&product).Error; err != nil {
				return err
			}

			for i := range product.Sizes {
				product.Sizes[i].ProductID = product.DBID
			}

			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&product.Sizes).Error; err != nil {
				return err
			}
		}

		return nil
	})

}
