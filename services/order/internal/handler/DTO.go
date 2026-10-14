package handler

type CreateOrderReq struct {
	UserID   int64     `json:"user_id"`
	Products []Product `json:"product"`
}

type Product struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"product_quantity"`
}

type CreateOrderResp struct {
	ID         int64   `json:"id"`
	UserID     int64   `json:"user_id"`
	ProductsID []int64 `json:"products_id"`
}
