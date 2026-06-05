package menu

import "time"

type MenuItem struct {
	ID          string    `json:"id" firestore:"id"`
	Name        string    `json:"name" firestore:"name"`
	Category    string    `json:"category" firestore:"category"`
	Price       float64   `json:"price" firestore:"price"`
	Unit        string    `json:"unit,omitempty" firestore:"unit"`
	Image       string    `json:"image,omitempty" firestore:"image"`
	Description string    `json:"description,omitempty" firestore:"description"`
	Active      bool      `json:"active" firestore:"active"`
	Order       int       `json:"order" firestore:"order"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" firestore:"updatedAt"`
}
