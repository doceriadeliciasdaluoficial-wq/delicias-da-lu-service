package order

type Order struct {
	Id         string      `json:"id"`
	UserId     string      `json:"user_id"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`

	ShippingAddress string `json:"shipping_address"`

	PaymentMethod string `json:"payment_method"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type OrderItem struct {
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
