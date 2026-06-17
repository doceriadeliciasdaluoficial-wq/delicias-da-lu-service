package home

import (
	"context"

	"delicias-da-lu-service.com/mod/internal/entity/config"
	"delicias-da-lu-service.com/mod/internal/repository/home"
)

type HomeUseCase interface {
	GetFeaturedCakes(ctx context.Context) ([]config.FeaturedCake, error)
	GetFeaturedCakeByID(ctx context.Context, id string) (*config.FeaturedCake, error)
	CreateFeaturedCake(ctx context.Context, cake *config.FeaturedCake) (*config.FeaturedCake, error)
	UpdateFeaturedCake(ctx context.Context, id string, cake *config.FeaturedCake) (*config.FeaturedCake, error)
	DeleteFeaturedCake(ctx context.Context, id string) error
	UpdateOrder(ctx context.Context, id string, order int) (*config.FeaturedCake, error)
}

type homeUseCaseImpl struct {
	homeRepository home.HomeRepository
}

func NewHomeUseCase(homeRepository home.HomeRepository) HomeUseCase {
	return homeUseCaseImpl{
		homeRepository: homeRepository,
	}
}

func (uc homeUseCaseImpl) GetFeaturedCakes(ctx context.Context) ([]config.FeaturedCake, error) {
	return uc.homeRepository.GetFeaturedCakes(ctx)
}

func (uc homeUseCaseImpl) GetFeaturedCakeByID(ctx context.Context, id string) (*config.FeaturedCake, error) {
	return uc.homeRepository.GetFeaturedCakeByID(ctx, id)
}

func (uc homeUseCaseImpl) CreateFeaturedCake(ctx context.Context, cake *config.FeaturedCake) (*config.FeaturedCake, error) {
	return uc.homeRepository.CreateFeaturedCake(ctx, cake)
}

func (uc homeUseCaseImpl) UpdateFeaturedCake(ctx context.Context, id string, cake *config.FeaturedCake) (*config.FeaturedCake, error) {
	return uc.homeRepository.UpdateFeaturedCake(ctx, id, cake)
}

func (uc homeUseCaseImpl) DeleteFeaturedCake(ctx context.Context, id string) error {
	return uc.homeRepository.DeleteFeaturedCake(ctx, id)
}

func (uc homeUseCaseImpl) UpdateOrder(ctx context.Context, id string, order int) (*config.FeaturedCake, error) {
	return uc.homeRepository.UpdateOrder(ctx, id, order)
}
