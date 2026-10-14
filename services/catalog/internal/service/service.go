package service

import (
	"context"
	"fmt"
	"shopik/services/catalog/internal/domain"
)

type CatalogRepository interface {
	GetByIDs(ctx context.Context, ids []int64) ([]domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
}

type CatalogService struct {
	repo CatalogRepository
}

func NewCatalogService(repo CatalogRepository) *CatalogService {
	return &CatalogService{
		repo: repo,
	}
}

func (s *CatalogService) GetProductsByIDs(ctx context.Context, ids []int64) ([]domain.Product, []int64, error) {
	var missingIDs []int64
	products, err := s.repo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, nil, err
	}
	foundIDs := make(map[int64]struct{}, len(products))
	for _, p := range products {
		foundIDs[p.ID] = struct{}{}
	}
	for _, id := range ids {
		_, ok := foundIDs[id]
		if !ok {
			missingIDs = append(missingIDs, id)
		}
	}
	return products, missingIDs, nil
}

func (s *CatalogService) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("Get all error: %w", err)
	}
	return products, nil
}
