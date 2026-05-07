package service

import (
	"context"
	"encoding/json"
	"wb_scraper/internal/domain"
)

type ProductsRepo interface {
	SaveBatch(ctx context.Context, products []domain.Product) error
}

type ProductsService struct {
	repo ProductsRepo
}

func NewProductService(repo ProductsRepo) *ProductsService {
	return &ProductsService{
		repo: repo,
	}
}

func (s *ProductsService) SyncProducts(ctx context.Context, data string) error {
	var wrapper struct {
		Products []domain.Product `json:"products"`
	}

	if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
		return err
	}

	for i := range wrapper.Products {
		for j := range wrapper.Products[i].Sizes {
			size := &wrapper.Products[i].Sizes[j]

			size.PriceBasic = float64(size.Price.Basic) / 100.0
			size.PriceProduct = float64(size.Price.Product) / 100.0
		}
	}

	return s.repo.SaveBatch(ctx, wrapper.Products)
}
