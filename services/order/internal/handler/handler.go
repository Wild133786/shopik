package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"shopik/services/order/internal/domain"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID int64, items []domain.OrderItem) (domain.Order, error)
}

type OrderHandler struct {
	srv OrderService
}

func NewOrderHandler(srv OrderService) *OrderHandler {
	return &OrderHandler{
		srv: srv,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq CreateOrderReq

	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var items []domain.OrderItem
	for _, item := range orderReq.Products {
		orderItem := domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		items = append(items, orderItem)
	}

	orderResp, err := h.srv.CreateOrder(r.Context(), orderReq.UserID, items)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			http.Error(w, "items, not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrItemForOrderOutOfStock) {
			http.Error(w, "items out of stock", http.StatusConflict)
			return
		}
		http.Error(w, "create order error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderResp)
}
