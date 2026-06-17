package cakebuilder

import (
	"context"
	"fmt"

	"delicias-da-lu-service.com/mod/internal/entity/cakebuilder"
	cakeBuilderRepo "delicias-da-lu-service.com/mod/internal/repository/cakebuilder"
)

type CakeBuilderUseCase interface {
	GetByType(ctx context.Context, componentType string, active *bool) ([]cakebuilder.CakeBuilderComponent, error)
	GetByID(ctx context.Context, componentType, id string) (*cakebuilder.CakeBuilderComponent, error)
	Create(ctx context.Context, componentType string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error)
	Update(ctx context.Context, componentType, id string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error)
	Delete(ctx context.Context, componentType, id string) error
	GetAll(ctx context.Context, active *bool) (map[string][]cakebuilder.CakeBuilderComponent, error)
	UpdateOrder(ctx context.Context, componentType, id string, order int) (*cakebuilder.CakeBuilderComponent, error)
}

type cakeBuilderUseCaseImpl struct {
	repository cakeBuilderRepo.CakeBuilderRepository
}

func NewCakeBuilderUseCase(repository cakeBuilderRepo.CakeBuilderRepository) CakeBuilderUseCase {
	return cakeBuilderUseCaseImpl{
		repository: repository,
	}
}

func (c cakeBuilderUseCaseImpl) GetByType(ctx context.Context, componentType string, active *bool) ([]cakebuilder.CakeBuilderComponent, error) {
	return c.repository.GetByType(ctx, componentType, active)
}

func (c cakeBuilderUseCaseImpl) GetByID(ctx context.Context, componentType, id string) (*cakebuilder.CakeBuilderComponent, error) {
	if componentType == "" {
		return nil, fmt.Errorf("component type is required")
	}
	return c.repository.GetByID(ctx, componentType, id)
}

func (c cakeBuilderUseCaseImpl) Create(ctx context.Context, componentType string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error) {
	if componentType == "" {
		return nil, fmt.Errorf("component type is required")
	}
	return c.repository.Create(ctx, componentType, component)
}

func (c cakeBuilderUseCaseImpl) Update(ctx context.Context, componentType, id string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error) {
	if componentType == "" {
		return nil, fmt.Errorf("component type is required")
	}
	return c.repository.Update(ctx, componentType, id, component)
}

func (c cakeBuilderUseCaseImpl) Delete(ctx context.Context, componentType, id string) error {
	if componentType == "" {
		return fmt.Errorf("component type is required")
	}
	return c.repository.Delete(ctx, componentType, id)
}

func (c cakeBuilderUseCaseImpl) GetAll(ctx context.Context, active *bool) (map[string][]cakebuilder.CakeBuilderComponent, error) {
	return c.repository.GetAll(ctx, active)
}

func (c cakeBuilderUseCaseImpl) UpdateOrder(ctx context.Context, componentType, id string, order int) (*cakebuilder.CakeBuilderComponent, error) {
	if componentType == "" {
		return nil, fmt.Errorf("component type is required")
	}
	return c.repository.UpdateOrder(ctx, componentType, id, order)
}
