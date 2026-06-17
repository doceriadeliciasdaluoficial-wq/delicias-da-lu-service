package menu

import (
	"context"
	"fmt"

	"delicias-da-lu-service.com/mod/internal/entity/menu"
	menuRepo "delicias-da-lu-service.com/mod/internal/repository/menu"
)

type MenuUseCase interface {
	GetAll(ctx context.Context, active *bool, category string) ([]menu.MenuItem, error)
	GetByID(ctx context.Context, categoryID, itemID string) (*menu.MenuItem, error)
	Create(ctx context.Context, categoryID string, item *menu.MenuItem) (*menu.MenuItem, error)
	Update(ctx context.Context, categoryID, itemID string, item *menu.MenuItem) (*menu.MenuItem, error)
	Delete(ctx context.Context, categoryID, itemID string) error
	UpdateOrder(ctx context.Context, categoryID, itemID string, order int) (*menu.MenuItem, error)
}

type menuUseCaseImpl struct {
	repository menuRepo.MenuRepository
}

func NewMenuUseCase(repository menuRepo.MenuRepository) MenuUseCase {
	return menuUseCaseImpl{
		repository: repository,
	}
}

func (m menuUseCaseImpl) GetAll(ctx context.Context, active *bool, category string) ([]menu.MenuItem, error) {
	return m.repository.GetAll(ctx, active, category)
}

func (m menuUseCaseImpl) GetByID(ctx context.Context, categoryID, itemID string) (*menu.MenuItem, error) {
	if categoryID == "" {
		return nil, fmt.Errorf("category ID is required")
	}
	return m.repository.GetByID(ctx, categoryID, itemID)
}

func (m menuUseCaseImpl) Create(ctx context.Context, categoryID string, item *menu.MenuItem) (*menu.MenuItem, error) {
	if categoryID == "" {
		return nil, fmt.Errorf("category ID is required")
	}
	return m.repository.Create(ctx, categoryID, item)
}

func (m menuUseCaseImpl) Update(ctx context.Context, categoryID, itemID string, item *menu.MenuItem) (*menu.MenuItem, error) {
	if categoryID == "" {
		return nil, fmt.Errorf("category ID is required")
	}
	return m.repository.Update(ctx, categoryID, itemID, item)
}

func (m menuUseCaseImpl) Delete(ctx context.Context, categoryID, itemID string) error {
	if categoryID == "" {
		return fmt.Errorf("category ID is required")
	}
	return m.repository.Delete(ctx, categoryID, itemID)
}

func (m menuUseCaseImpl) UpdateOrder(ctx context.Context, categoryID, itemID string, order int) (*menu.MenuItem, error) {
	if categoryID == "" {
		return nil, fmt.Errorf("category ID is required")
	}
	return m.repository.UpdateOrder(ctx, categoryID, itemID, order)
}
