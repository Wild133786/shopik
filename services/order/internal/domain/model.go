package domain

import "time"

type OrderStatus int

const (
	StatusCreated OrderStatus = iota
	StatusPaid
	StatusShipped
	StatusCancelled
)

type Order struct {
	ID         int64
	UserID     int64
	Status     OrderStatus
	Items      []OrderItem
	TotalPrice int64
	CreatedAt  time.Time
}

type OrderItem struct {
	ProductID int64
	Quantity  int64
}

type Product struct {
	ID    int64
	Name  string
	Price int64
	Stock int64
}
