package domain

import "errors"

var (
	ErrItemNotFound           = errors.New("item not found in catalog")
	ErrItemForOrderOutOfStock = errors.New("some items from order, out of stock")
)
