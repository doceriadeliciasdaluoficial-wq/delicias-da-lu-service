package cakebuilder

import "time"

type CakeBuilderComponent struct {
	ID              string    `json:"id" firestore:"id"`
	Name            string    `json:"name" firestore:"name"`
	Label           string    `json:"label" firestore:"label"`
	Type            string    `json:"type" firestore:"type"`
	Price           float64   `json:"price" firestore:"price"`
	Value           float64   `json:"value" firestore:"value"`
	Description     string    `json:"description,omitempty" firestore:"description"`
	FullDescription *string   `json:"fullDescription,omitempty" firestore:"fullDescription"`
	Weight          *string   `json:"weight,omitempty" firestore:"weight"`
	Servings        *string   `json:"servings,omitempty" firestore:"servings"`
	Note            *string   `json:"note,omitempty" firestore:"note"`
	Image           string    `json:"image,omitempty" firestore:"image"`
	Active          bool      `json:"active" firestore:"active"`
	Order           int       `json:"order" firestore:"order"`
	CreatedAt       time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt" firestore:"updatedAt"`
}
