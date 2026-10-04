package repository

import (
	"cmp"
	"context"
	"shopik/services/catalog/internal/domain"
	"slices"
)

type InMemoryRepository struct {
	products (map[int64]domain.Product)
}

func NewInMemoryRepository() *InMemoryRepository {
	r := &InMemoryRepository{products: make(map[int64]domain.Product)}
	r.seed()
	return r
}

func (r *InMemoryRepository) seed() {
	r.products[1] = domain.Product{
		ID:    1,
		Name:  "Gamepad Logitech",
		Price: 600000,
		Stock: 6,
	}
	r.products[2] = domain.Product{
		ID:    2,
		Name:  "Headset Logitech",
		Price: 1200000,
		Stock: 9,
	}
	r.products[3] = domain.Product{
		ID:    3,
		Name:  "Mouse Logitech",
		Price: 400000,
		Stock: 3,
	}
}

func (r *InMemoryRepository) GetByIDs(ctx context.Context, ids []int64) ([]domain.Product, error) {
	var result []domain.Product
	for _, id := range ids {
		product, ok := r.products[id]
		if !ok {
			continue
		}
		result = append(result, product)
	}
	return result, nil
}

func (r *InMemoryRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	var result []domain.Product
	for _, product := range r.products {
		result = append(result, product)
	}
	slices.SortFunc(result, func(a, b domain.Product) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return result, nil
}
