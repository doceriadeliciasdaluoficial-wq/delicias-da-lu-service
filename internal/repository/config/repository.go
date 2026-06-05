package config

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"delicias-da-lu-service.com/mod/internal/entity/config"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
)

const configDocID = "site-config"

type ConfigRepository interface {
	Get(ctx context.Context) (*config.SiteConfig, error)
	Update(ctx context.Context, cfg *config.SiteConfig) (*config.SiteConfig, error)
}

type configRepositoryImpl struct {
	client *firestore.Client
}

func NewConfigRepository(client *firestore.Client) ConfigRepository {
	return configRepositoryImpl{
		client: client,
	}
}

func (r configRepositoryImpl) Get(ctx context.Context) (*config.SiteConfig, error) {
	doc, err := r.client.Collection("config").Doc(configDocID).Get(ctx)
	if err != nil {
		return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
			Title:      "Configuration Not Found",
			Detail:     "Unable to retrieve site configuration",
			HTTPStatus: http.StatusInternalServerError,
			Instance:   "localhost:8080/v1/config",
			Severity:   problemdetails.Err,
		})
	}

	var cfg config.SiteConfig
	if err := doc.DataTo(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (r configRepositoryImpl) Update(ctx context.Context, cfg *config.SiteConfig) (*config.SiteConfig, error) {
	_, err := r.client.Collection("config").Doc(configDocID).Set(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
