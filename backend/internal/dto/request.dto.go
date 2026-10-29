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
