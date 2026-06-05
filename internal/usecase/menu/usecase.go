package menu

import (
	"context"

	"delicias-da-lu-service.com/mod/internal/entity/menu"
	menuRepo "delicias-da-lu-service.com/mod/internal/repository/menu"
)

type MenuUseCase interface {
	GetAll(ctx context.Context, active *bool, category string) ([]menu.MenuItem, error)
	GetByID(ctx context.Context, id string) (*menu.MenuItem, error)
	Create(ctx context.Context, item *menu.MenuItem) (*menu.MenuItem, error)
	Update(ctx context.Context, id string, item *menu.MenuItem) (*menu.MenuItem, error)
	Delete(ctx context.Context, id string) error
	UpdateOrder(ctx context.Context, id string, order int) (*menu.MenuItem, error)
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

func (m menuUseCaseImpl) GetByID(ctx context.Context, id string) (*menu.MenuItem, error) {
	return m.repository.GetByID(ctx, id)
}

func (m menuUseCaseImpl) Create(ctx context.Context, item *menu.MenuItem) (*menu.MenuItem, error) {
	return m.repository.Create(ctx, item)
}

func (m menuUseCaseImpl) Update(ctx context.Context, id string, item *menu.MenuItem) (*menu.MenuItem, error) {
	return m.repository.Update(ctx, id, item)
}

func (m menuUseCaseImpl) Delete(ctx context.Context, id string) error {
	return m.repository.Delete(ctx, id)
}

func (m menuUseCaseImpl) UpdateOrder(ctx context.Context, id string, order int) (*menu.MenuItem, error) {
	return m.repository.UpdateOrder(ctx, id, order)
}
