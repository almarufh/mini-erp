package dto

type CreateProductReq struct {
	SKU      string `json:"sku" binding:"required" example:"SKU-KB-001"`
	Name     string `json:"name" binding:"required" example:"Mechanical Keyboard Wireless Red Switch"`
	Price    int    `json:"price" binding:"gte=0" example:"850000"`
	Stock    int    `json:"stock" binding:"gte=0" example:"25"`
	IsActive *bool  `json:"is_active" example:"true"`
}

type UpdateProductReq struct {
	Name     *string `json:"name,omitempty" example:"Mechanical Keyboard Wireless RGB"`
	Price    *int    `json:"price,omitempty" example:"950000"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}

type AdjustStockReq struct {
	Type string  `json:"type" binding:"required,oneof=in out" example:"out"`
	Qty  int     `json:"qty" binding:"required,gt=0" example:"5"`
	Note *string `json:"note,omitempty" example:"Pengurangan stok manual"`
}

type CreateOrderItemReq struct {
	SKU string `json:"sku" binding:"required" example:"SKU-KB-001"`
	Qty int    `json:"qty" binding:"required,gt=0" example:"2"`
}

type CreateOrderReq struct {
	Customer string               `json:"customer" binding:"required" example:"Budi Santoso"`
	Items    []CreateOrderItemReq `json:"items" binding:"required,gt=0,dive"`
}
