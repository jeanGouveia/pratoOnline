package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID         uint
	Status     OrderStatus
	TotalPrice float64
	Notes      string
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Preenchido sob demanda
	Items []OrderItem
}
