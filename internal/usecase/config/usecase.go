package config

import (
	"context"

	"delicias-da-lu-service.com/mod/internal/entity/config"
	configRepo "delicias-da-lu-service.com/mod/internal/repository/config"
)

type ConfigUseCase interface {
	Get(ctx context.Context) (*config.SiteConfig, error)
	Update(ctx context.Context, cfg *config.SiteConfig) (*config.SiteConfig, error)
}

type configUseCaseImpl struct {
	repository configRepo.ConfigRepository
}

func NewConfigUseCase(repository configRepo.ConfigRepository) ConfigUseCase {
	return configUseCaseImpl{
		repository: repository,
	}
}

func (c configUseCaseImpl) Get(ctx context.Context) (*config.SiteConfig, error) {
	return c.repository.Get(ctx)
}

func (c configUseCaseImpl) Update(ctx context.Context, cfg *config.SiteConfig) (*config.SiteConfig, error) {
	return c.repository.Update(ctx, cfg)
}
