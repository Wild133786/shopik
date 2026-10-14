package service

import (
	"context"
	"fmt"
	"shopik/services/order/internal/domain"
	"time"
)

type OrderRepository interface {
	Save(ctx context.Context, order domain.Order) (int64, error)
}

type CatalogProduct struct {
	ID    int64
	Name  string
	Price int64
	Stock int64
}

type CatalogClient interface {
	GetProductsByIDs(ctx context.Context, ids []int64) ([]CatalogProduct, []int64, error)
}

type OrderService struct {
	repo    OrderRepository
	catalog CatalogClient
}

func NewOrderService(repo OrderRepository, catalog CatalogClient) *OrderService {
	return &OrderService{
		repo:    repo,
		catalog: catalog,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID int64, items []domain.OrderItem) (domain.Order, error) {
	var itemsIDs []int64
	for _, item := range items {
		itemsIDs = append(itemsIDs, item.ProductID)
	}
	productsCatalogResp, missingIDs, err := s.catalog.GetProductsByIDs(ctx, itemsIDs)
	if err != nil {
		return domain.Order{}, err
	}
	if len(missingIDs) != 0 {
		return domain.Order{}, domain.ErrItemNotFound
	}
	var outOfStocksItems []domain.OrderItem
	var orderTotalPrice int64
	mapProductsCatalogResp := make(map[int64]CatalogProduct, len(productsCatalogResp))
	for _, productCatalogResp := range productsCatalogResp {
		mapProductsCatalogResp[productCatalogResp.ID] = productCatalogResp
	}
	for _, item := range items {
		orderTotalPrice += mapProductsCatalogResp[item.ProductID].Price * item.Quantity
		if mapProductsCatalogResp[item.ProductID].Stock < item.Quantity {
			outOfStocksItems = append(outOfStocksItems, item)
		}
	}
	if len(outOfStocksItems) != 0 {
		return domain.Order{}, fmt.Errorf("%w: %+v", domain.ErrItemForOrderOutOfStock, outOfStocksItems)
	}
	order := domain.Order{
		UserID:     userID,
		Items:      items,
		Status:     domain.StatusCreated,
		CreatedAt:  time.Now(),
		TotalPrice: orderTotalPrice,
	}

	order.ID, err = s.repo.Save(ctx, order)
	if err != nil {
		return domain.Order{}, fmt.Errorf("saving order error: %w", err) // TODO domain error
	}
	return order, nil
}
