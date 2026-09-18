package handler

import (
	"encoding/json"
	"net/http"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var products []ProductForResponse
	products = []ProductForResponse{
		{
			ID:    1,
			Name:  "Keyboard Ajazz",
			Price: 599000,
		},
		{
			ID:    2,
			Name:  "Mouse Razer",
			Price: 460000,
		},
		{
			ID:    3,
			Name:  "Headset HyperX",
			Price: 720000,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
