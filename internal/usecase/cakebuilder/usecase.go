package cakebuilder

import (
	"context"

	"delicias-da-lu-service.com/mod/internal/entity/cakebuilder"
	cakebuilderRepo "delicias-da-lu-service.com/mod/internal/repository/cakebuilder"
)

type CakeBuilderUseCase interface {
	GetByType(ctx context.Context, componentType string, active *bool) ([]cakebuilder.CakeBuilderComponent, error)
	GetByID(ctx context.Context, componentType, id string) (*cakebuilder.CakeBuilderComponent, error)
	Create(ctx context.Context, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error)
	Update(ctx context.Context, id string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error)
	Delete(ctx context.Context, componentType, id string) error
	GetAll(ctx context.Context, active *bool) (map[string][]cakebuilder.CakeBuilderComponent, error)
}

type cakeBuilderUseCaseImpl struct {
	repository cakebuilderRepo.CakeBuilderRepository
}

func NewCakeBuilderUseCase(repository cakebuilderRepo.CakeBuilderRepository) CakeBuilderUseCase {
	return cakeBuilderUseCaseImpl{
		repository: repository,
	}
}

func (c cakeBuilderUseCaseImpl) GetByType(ctx context.Context, componentType string, active *bool) ([]cakebuilder.CakeBuilderComponent, error) {
	return c.repository.GetByType(ctx, componentType, active)
}

func (c cakeBuilderUseCaseImpl) GetByID(ctx context.Context, componentType, id string) (*cakebuilder.CakeBuilderComponent, error) {
	return c.repository.GetByID(ctx, componentType, id)
}

func (c cakeBuilderUseCaseImpl) Create(ctx context.Context, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error) {
	return c.repository.Create(ctx, component)
}

func (c cakeBuilderUseCaseImpl) Update(ctx context.Context, id string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error) {
	return c.repository.Update(ctx, id, component)
}

func (c cakeBuilderUseCaseImpl) Delete(ctx context.Context, componentType, id string) error {
	return c.repository.Delete(ctx, componentType, id)
}

func (c cakeBuilderUseCaseImpl) GetAll(ctx context.Context, active *bool) (map[string][]cakebuilder.CakeBuilderComponent, error) {
	return c.repository.GetAll(ctx, active)
}
