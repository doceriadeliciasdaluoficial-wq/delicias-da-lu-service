package cakebuilder

import "time"

type CakeBuilderComponent struct {
	ID        string    `json:"id" firestore:"id"`
	Name      string    `json:"name" firestore:"name"`
	Type      string    `json:"type" firestore:"type"` // massa, recheio, cobertura, decoracao
	Price     float64   `json:"price" firestore:"price"`
	Image     string    `json:"image,omitempty" firestore:"image"`
	Active    bool      `json:"active" firestore:"active"`
	Order     int       `json:"order" firestore:"order"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt"`
}
