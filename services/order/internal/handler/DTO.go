package handler

type CreateOrderReq struct {
	UserID     int64   `json:"user_id"`
	ProductsID []int64 `json:"products_id"`
}

type CreateOrderResp struct {
	ID         int64   `json:"id"`
	UserID     int64   `json:"user_id"`
	ProductsID []int64 `json:"products_id"`
}
