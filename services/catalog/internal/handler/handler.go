package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"shopik/services/catalog/internal/domain"
)

type Service interface {
	GetProductsByIDs(ctx context.Context, ids []int64) ([]domain.Product, []int64, error)
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
}

type CatalogHandler struct {
	service Service
}

func NewCatalogHandler(service Service) *CatalogHandler {
	return &CatalogHandler{service: service}
}

func (h *CatalogHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var respProducts []ProductForResponse
	for _, product := range products {
		respProduct := ProductForResponse{
			ID:    product.ID,
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Stock,
		}
		respProducts = append(respProducts, respProduct)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respProducts)
}
