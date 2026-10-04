package repository

import (
	"cmp"
	"context"
	"shopik/services/order/internal/domain"
	"slices"
	"sync"
	"time"
)

type InMemoryRepository struct {
	orders (map[int64]domain.Order)
	mu     sync.Mutex
	nextID int64
}

func NewInMemoryRepository() *InMemoryRepository {
	r := &InMemoryRepository{
		orders: make(map[int64]domain.Order),
		nextID: 1,
	}
	r.seed()
	return r
}

func (r *InMemoryRepository) Save(ctx context.Context, order domain.Order) (int64, error) { // TODO id generate, sync
	r.mu.Lock()
	defer r.mu.Unlock()
	order.ID = r.nextID
	r.orders[order.ID] = order
	r.nextID++
	return order.ID, nil
}

func (r *InMemoryRepository) GetAll(ctx context.Context) ([]domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var respOrders []domain.Order
	for _, order := range r.orders {
		respOrders = append(respOrders, order)
	}
	slices.SortFunc(respOrders, func(a, b domain.Order) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return respOrders, nil
}

func (r *InMemoryRepository) seed() {
	r.orders[0] = domain.Order{
		ID:     0,
		UserID: 5,
		Status: domain.StatusCreated,
		Items: []domain.OrderItem{
			{
				ProductID: 1,
				Quantity:  124,
			},
			{
				ProductID: 15,
				Quantity:  1,
			},
		},
		TotalPrice: 2142,
		CreatedAt:  time.Now(),
	}
}
