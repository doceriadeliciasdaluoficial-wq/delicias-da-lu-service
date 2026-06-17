package cakebuilder

import (
	"testing"
	"time"

	"delicias-da-lu-service.com/mod/internal/entity/cakebuilder"
	"github.com/stretchr/testify/assert"
)

func createTestComponent(id, name, cType string, price float64) *cakebuilder.CakeBuilderComponent {
	return &cakebuilder.CakeBuilderComponent{
		ID:        id,
		Name:      name,
		Type:      cType,
		Price:     price,
		Image:     "test.jpg",
		Active:    true,
		Order:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestCakeBuilderComponentCreation(t *testing.T) {
	component := createTestComponent("massa-001", "Massa Branca", "massas", 15.00)

	assert.Equal(t, "massa-001", component.ID)
	assert.Equal(t, "Massa Branca", component.Name)
	assert.Equal(t, "massas", component.Type)
	assert.Equal(t, 15.00, component.Price)
	assert.True(t, component.Active)
}

func TestCakeBuilderComponentsByType(t *testing.T) {
	components := map[string][]*cakebuilder.CakeBuilderComponent{
		"massas": {
			createTestComponent("massa-001", "Massa Branca", "massas", 15.00),
			createTestComponent("massa-002", "Massa Chocolate", "massas", 15.00),
		},
		"recheios": {
			createTestComponent("recheio-001", "Morango", "recheios", 10.00),
			createTestComponent("recheio-002", "Chocolate", "recheios", 10.00),
		},
		"coberturas": {
			createTestComponent("cobertura-001", "Ganache", "coberturas", 8.00),
		},
	}

	assert.Len(t, components, 3)
	assert.Len(t, components["massas"], 2)
	assert.Len(t, components["recheios"], 2)
	assert.Len(t, components["coberturas"], 1)

	for typeKey, comps := range components {
		for _, comp := range comps {
			assert.Equal(t, typeKey, comp.Type)
		}
	}
}

func TestCakeBuilderValidTypes(t *testing.T) {
	validTypes := []string{"massas", "recheios", "coberturas", "decoracoes", "sizes"}

	for _, vType := range validTypes {
		t.Run(vType, func(t *testing.T) {
			comp := createTestComponent("test-001", "Test", vType, 10.00)
			assert.Equal(t, vType, comp.Type)
		})
	}
}

func TestCakeBuilderComponentOrdering(t *testing.T) {
	components := []*cakebuilder.CakeBuilderComponent{
		createTestComponent("massa-001", "Massa Branca", "massas", 15.00),
		createTestComponent("massa-002", "Massa Chocolate", "massas", 15.00),
		createTestComponent("massa-003", "Massa Red Velvet", "massas", 18.00),
	}

	for i, comp := range components {
		comp.Order = i + 1
	}

	for i, comp := range components {
		assert.Equal(t, i+1, comp.Order)
	}
}

func TestCakeBuilderComponentPricing(t *testing.T) {
	tests := []struct {
		name  string
		cType string
		price float64
	}{
		{"Massa", "massas", 15.00},
		{"Recheio", "recheios", 10.50},
		{"Cobertura", "coberturas", 8.99},
		{"Decoração", "decoracoes", 5.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := createTestComponent("test-001", tt.name, tt.cType, tt.price)
			assert.Equal(t, tt.price, comp.Price)
		})
	}
}

func TestCakeBuilderOptionalFields(t *testing.T) {
	comp := createTestComponent("massa-001", "Massa Branca", "massas", 15.00)

	description := "A fine white cake"
	comp.FullDescription = &description

	weight := "500g"
	comp.Weight = &weight

	servings := "10 pessoas"
	comp.Servings = &servings

	assert.NotNil(t, comp.FullDescription)
	assert.Equal(t, "A fine white cake", *comp.FullDescription)

	assert.NotNil(t, comp.Weight)
	assert.Equal(t, "500g", *comp.Weight)

	assert.NotNil(t, comp.Servings)
	assert.Equal(t, "10 pessoas", *comp.Servings)
}

func TestCakeBuilderComponentUpdate(t *testing.T) {
	comp := createTestComponent("massa-001", "Massa Branca", "massas", 15.00)

	originalCreatedAt := comp.CreatedAt
	comp.UpdatedAt = time.Now().Add(time.Hour)
	comp.Price = 18.00

	assert.Equal(t, originalCreatedAt, comp.CreatedAt)
	assert.NotEqual(t, originalCreatedAt, comp.UpdatedAt)
	assert.Equal(t, 18.00, comp.Price)
}
