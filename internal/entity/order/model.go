package order

import "time"

type Order struct {
	ID           string       `json:"id" firestore:"id"`
	Items        []OrderItem  `json:"items" firestore:"items"`
	CustomerInfo CustomerInfo `json:"customerInfo" firestore:"customerInfo"`
	Status       string       `json:"status" firestore:"status"` // pending, confirmed, preparing, ready, delivered, cancelled
	TotalPrice   float64      `json:"totalPrice" firestore:"totalPrice"`
	CreatedAt    time.Time    `json:"createdAt" firestore:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt" firestore:"updatedAt"`
}

type OrderItem struct {
	Type              string                 `json:"type" firestore:"type"` // menu or cakeBuilder
	MenuItemID        string                 `json:"menuItemId,omitempty" firestore:"menuItemId"`
	CakeCustomization map[string]interface{} `json:"cakeCustomization,omitempty" firestore:"cakeCustomization"`
	Quantity          int                    `json:"quantity" firestore:"quantity"`
	UnitPrice         float64                `json:"unitPrice" firestore:"unitPrice"`
	Subtotal          float64                `json:"subtotal" firestore:"subtotal"`
}

type CustomerInfo struct {
	Name         string `json:"name" firestore:"name"`
	Phone        string `json:"phone" firestore:"phone"`
	Email        string `json:"email,omitempty" firestore:"email"`
	DeliveryDate string `json:"deliveryDate,omitempty" firestore:"deliveryDate"` // date format
	Notes        string `json:"notes,omitempty" firestore:"notes"`
}

type OrderListResponse struct {
	Total  int     `json:"total"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
	Data   []Order `json:"data"`
}
