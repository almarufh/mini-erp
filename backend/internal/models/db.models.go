package models

import "time"

type Product struct {
	ID             int        `db:"id"`
	SKU            string     `db:"sku"`
	Name           string     `db:"name"`
	Price          int        `db:"price"`
	Stock          int        `db:"stock"`
	IsActive       bool       `db:"is_active"`
	CreatedAt      time.Time  `db:"created_at"`
	LastMovementAt *time.Time `db:"last_movement_at"`
}

type StockMovement struct {
	ID        int       `db:"id"`
	ProductID int       `db:"product_id"`
	Type      string    `db:"type"`
	Qty       int       `db:"qty"`
	Note      *string   `db:"note"`
	CreatedAt time.Time `db:"created_at"`
}

type SalesOrder struct {
	ID          int        `db:"id"`
	Customer    string     `db:"customer"`
	Status      string     `db:"status"`
	Total       int        `db:"total"`
	CreatedAt   time.Time  `db:"created_at"`
	CancelledAt *time.Time `db:"cancelled_at"`
}

type SalesOrderItem struct {
	ID        int    `db:"id"`
	OrderID   int    `db:"order_id"`
	SKU       string `db:"sku"`
	Qty       int    `db:"qty"`
	UnitPrice int    `db:"unit_price"`
}
