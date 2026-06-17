package config

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"delicias-da-lu-service.com/mod/internal/entity/config"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		if status.Code(err) == codes.NotFound {
			return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
				Title:      "Configuration Not Found",
				Detail:     "Site configuration has not been initialized. Please set up the configuration in the admin panel.",
				HTTPStatus: http.StatusNotFound,
				Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/config/public",
				Severity:   problemdetails.Err,
			})
		}
		return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/database-error",
			Title:      "Database Error",
			Detail:     "Failed to retrieve configuration: " + err.Error(),
			HTTPStatus: http.StatusInternalServerError,
			Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/config/public",
			Severity:   problemdetails.Err,
		})
	}

	var cfg config.SiteConfig
	if err := doc.DataTo(&cfg); err != nil {
		return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/parse-error",
			Title:      "Configuration Parse Error",
			Detail:     "Failed to parse configuration data: " + err.Error(),
			HTTPStatus: http.StatusInternalServerError,
			Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/config/public",
			Severity:   problemdetails.Err,
		})
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
