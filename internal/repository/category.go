package repository

import (
	"context"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"wb_scraper/internal/domain"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (repo *CategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	return repo.db.Table("wb_scraper.categories").WithContext(ctx).Create(category).Error
}

func (repo *CategoryRepository) Delete(ctx context.Context, id int) error {
	res := repo.db.Table("wb_scraper.categories").
		WithContext(ctx).
		Where("wb_id = ?", id).
		Delete(&domain.Category{})
	if res.Error == nil && res.RowsAffected == 0 {
		log.Printf("Warning: category with id [%d] not found (nothing to delete)", id)
	}
	return res.Error
}

func (repo *CategoryRepository) SaveBatch(ctx context.Context, categories []domain.Category) error {
	return repo.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "wb_id"}},
			UpdateAll: true,
		}).
		CreateInBatches(categories, 100).Error
}

func (repo *CategoryRepository) GetAll(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category

	err := repo.db.WithContext(ctx).Where("url <> ''").Find(&categories).Error
	return categories, err
}
