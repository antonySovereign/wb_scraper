package service

import (
	"context"
	"encoding/json"
	"net/http"

	"wb_scraper/internal/domain"
)

type CategoryRepo interface {
	SaveBatch(ctx context.Context, categories []domain.Category) error
}

type CategoryService struct {
	repo CategoryRepo
}

func NewCategoryService(repo CategoryRepo) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) SyncCategories(ctx context.Context, url string) error {
	// get data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var categories []domain.Category

	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return err
	}

	var flatCategories []domain.Category
	for _, cat := range categories {
		flatCategories = append(flatCategories, flatten(cat, 0)...)
	}

	return s.repo.SaveBatch(ctx, flatCategories)
}

func flatten(cat domain.Category, parentId int) []domain.Category {
	cat.Parent = parentId
	result := []domain.Category{cat}
	for _, child := range cat.Childs {
		result = append(result, flatten(child, cat.ID)...)
	}
	return result
}
