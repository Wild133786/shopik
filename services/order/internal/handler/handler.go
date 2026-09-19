package handler

import (
	"encoding/json"
	"net/http"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq CreateOrderReq

	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderResp := CreateOrderResp{
		ID:         1,
		UserID:     orderReq.UserID,
		ProductsID: orderReq.ProductsID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderResp)
}
