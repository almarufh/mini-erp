package dto

import "time"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Results any    `json:"results,omitempty"`
}

type PaginatedProductsResponse struct {
	Data      []ProductResponse `json:"data"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
	TotalData int               `json:"total_data"`
	TotalPage int               `json:"total_page"`
	NextPage  string            `json:"next_page"`
	PrevPage  string            `json:"prev_page"`
}

type ProductResponse struct {
	ID             int        `json:"id"`
	SKU            string     `json:"sku"`
	Name           string     `json:"name"`
	Price          int        `json:"price"`
	Stock          int        `json:"stock"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	LastMovementAt *time.Time `json:"last_movement_at"`
}

type ErrorEnvelope struct {
	Error struct {
		Code      string `json:"code"`
		SKU       string `json:"sku,omitempty"`
		Available *int   `json:"available,omitempty"`
	} `json:"error"`
}
