package menu

import (
	"testing"
	"time"

	"delicias-da-lu-service.com/mod/internal/entity/menu"
	"github.com/stretchr/testify/assert"
)

func createTestMenuItem(id, name, category string, price float64) *menu.MenuItem {
	return &menu.MenuItem{
		ID:        id,
		Name:      name,
		Category:  category,
		Price:     price,
		Unit:      "Un",
		Image:     "test.jpg",
		Active:    true,
		Order:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestMenuItemCreation(t *testing.T) {
	item := createTestMenuItem("bolos-001", "Bolo de Chocolate", "bolos", 45.00)

	assert.Equal(t, "bolos-001", item.ID)
	assert.Equal(t, "Bolo de Chocolate", item.Name)
	assert.Equal(t, "bolos", item.Category)
	assert.Equal(t, 45.00, item.Price)
	assert.True(t, item.Active)
}

func TestMenuItemsCategorization(t *testing.T) {
	items := map[string][]*menu.MenuItem{
		"bolos": {
			createTestMenuItem("bolos-001", "Bolo de Chocolate", "bolos", 45.00),
			createTestMenuItem("bolos-002", "Bolo de Morango", "bolos", 50.00),
		},
		"docesSimples": {
			createTestMenuItem("simples-001", "Brigadeiro", "docesSimples", 2.00),
		},
	}

	assert.Len(t, items, 2)
	assert.Len(t, items["bolos"], 2)
	assert.Len(t, items["docesSimples"], 1)

	for category, itemList := range items {
		for _, item := range itemList {
			assert.Equal(t, category, item.Category)
		}
	}
}

func TestMenuItemPricing(t *testing.T) {
	tests := []struct {
		name  string
		price float64
	}{
		{"Bolo Simples", 25.00},
		{"Bolo Premium", 45.00},
		{"Bolo Gourmet", 65.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := createTestMenuItem("test-"+tt.name, tt.name, "bolos", tt.price)
			assert.Equal(t, tt.price, item.Price)
		})
	}
}

func TestMenuItemOrdering(t *testing.T) {
	items := []*menu.MenuItem{
		createTestMenuItem("bolo-001", "Bolo Chocolate", "bolos", 45.00),
		createTestMenuItem("bolo-002", "Bolo Morango", "bolos", 50.00),
		createTestMenuItem("bolo-003", "Bolo Baunilha", "bolos", 45.00),
	}

	for i, item := range items {
		item.Order = i + 1
	}

	for i, item := range items {
		assert.Equal(t, i+1, item.Order)
	}
}

func TestMenuItemValidation(t *testing.T) {
	tests := []struct {
		name     string
		itemID   string
		itemName string
		isValid  bool
	}{
		{"Valid", "bolos-001", "Bolo", true},
		{"Empty ID", "", "Bolo", false},
		{"Valid name", "item-001", "Bolo de Chocolate", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isValid {
				item := createTestMenuItem(tt.itemID, tt.itemName, "bolos", 25.00)
				assert.NotNil(t, item)
				assert.Equal(t, tt.itemID, item.ID)
				assert.Equal(t, tt.itemName, item.Name)
			}
		})
	}
}

func TestMenuItemTimestamps(t *testing.T) {
	item := createTestMenuItem("test-001", "Test Item", "bolos", 25.00)

	assert.False(t, item.CreatedAt.IsZero())
	assert.False(t, item.UpdatedAt.IsZero())

	originalCreatedAt := item.CreatedAt
	item.UpdatedAt = time.Now().Add(time.Hour)

	assert.Equal(t, originalCreatedAt, item.CreatedAt)
	assert.NotEqual(t, originalCreatedAt, item.UpdatedAt)
}

func TestMenuItemActivity(t *testing.T) {
	activeItem := createTestMenuItem("active-001", "Active Item", "bolos", 25.00)
	activeItem.Active = true

	inactiveItem := createTestMenuItem("inactive-001", "Inactive Item", "bolos", 25.00)
	inactiveItem.Active = false

	assert.True(t, activeItem.Active)
	assert.False(t, inactiveItem.Active)
}
