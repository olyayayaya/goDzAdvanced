package order

type CreateOrderRequest struct {
	ProductIDs []uint `json:"product_ids" validate:"required,min=1"`
	Date       string `json:"date" validate:"required"`
}

type OrderResponse struct {
	ID        uint               `json:"id"`
	UserID    uint               `json:"user_id"`
	Date      string             `json:"date"`
	Products  []ProductResponse  `json:"products"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
}

type ProductResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}
