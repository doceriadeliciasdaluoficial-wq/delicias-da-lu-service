package home

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"delicias-da-lu-service.com/mod/internal/entity/config"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HomeRepository interface {
	GetFeaturedCakes(ctx context.Context) ([]config.FeaturedCake, error)
	GetFeaturedCakeByID(ctx context.Context, id string) (*config.FeaturedCake, error)
	CreateFeaturedCake(ctx context.Context, cake *config.FeaturedCake) (*config.FeaturedCake, error)
	UpdateFeaturedCake(ctx context.Context, id string, cake *config.FeaturedCake) (*config.FeaturedCake, error)
	DeleteFeaturedCake(ctx context.Context, id string) error
	UpdateOrder(ctx context.Context, id string, order int) (*config.FeaturedCake, error)
}

type homeRepositoryImpl struct {
	client *firestore.Client
}

const homeCollection = "home"
const featuredCakesSubcollection = "featuredCakes"

func NewHomeRepository(client *firestore.Client) HomeRepository {
	return &homeRepositoryImpl{client: client}
}

// GetFeaturedCakes retrieves all featured cakes ordered by order field
func (r *homeRepositoryImpl) GetFeaturedCakes(ctx context.Context) ([]config.FeaturedCake, error) {
	var cakes []config.FeaturedCake

	docs, err := r.client.Collection(homeCollection).Doc("default").Collection(featuredCakesSubcollection).
		OrderBy("order", firestore.Asc).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch featured cakes: %w", err)
	}

	for _, doc := range docs {
		var cake config.FeaturedCake
		if err := doc.DataTo(&cake); err != nil {
			return nil, err
		}
		cakes = append(cakes, cake)
	}

	return cakes, nil
}

// GetFeaturedCakeByID retrieves a specific featured cake
func (r *homeRepositoryImpl) GetFeaturedCakeByID(ctx context.Context, id string) (*config.FeaturedCake, error) {
	doc, err := r.client.Collection(homeCollection).Doc("default").Collection(featuredCakesSubcollection).Doc(id).Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
				Title:      "Featured Cake Not Found",
				Detail:     fmt.Sprintf("No featured cake found with ID: %s", id),
				HTTPStatus: http.StatusNotFound,
				Instance:   fmt.Sprintf("https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/home/featured-cakes/%s", id),
				Severity:   problemdetails.Err,
			})
		}
		return nil, err
	}

	var cake config.FeaturedCake
	if err := doc.DataTo(&cake); err != nil {
		return nil, err
	}

	return &cake, nil
}

// CreateFeaturedCake adds a new featured cake
func (r *homeRepositoryImpl) CreateFeaturedCake(ctx context.Context, cake *config.FeaturedCake) (*config.FeaturedCake, error) {
	// Ensure home document exists
	if _, err := r.client.Collection(homeCollection).Doc("default").Get(ctx); err != nil {
		if status.Code(err) == codes.NotFound {
			if _, err := r.client.Collection(homeCollection).Doc("default").Set(ctx, map[string]interface{}{
				"createdAt": time.Now(),
				"updatedAt": time.Now(),
			}); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	cake.CreatedAt = time.Now()
	cake.UpdatedAt = time.Now()

	_, err := r.client.Collection(homeCollection).Doc("default").Collection(featuredCakesSubcollection).Doc(cake.ID).Set(ctx, cake)
	if err != nil {
		return nil, err
	}

	return cake, nil
}

// UpdateFeaturedCake modifies an existing featured cake
func (r *homeRepositoryImpl) UpdateFeaturedCake(ctx context.Context, id string, cake *config.FeaturedCake) (*config.FeaturedCake, error) {
	existing, err := r.GetFeaturedCakeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	cake.ID = id
	cake.CreatedAt = existing.CreatedAt
	cake.UpdatedAt = time.Now()

	_, err = r.client.Collection(homeCollection).Doc("default").Collection(featuredCakesSubcollection).Doc(id).Set(ctx, cake)
	if err != nil {
		return nil, err
	}

	return cake, nil
}

// DeleteFeaturedCake removes a featured cake
func (r *homeRepositoryImpl) DeleteFeaturedCake(ctx context.Context, id string) error {
	_, err := r.GetFeaturedCakeByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = r.client.Collection(homeCollection).Doc("default").Collection(featuredCakesSubcollection).Doc(id).Delete(ctx)
	return err
}

// UpdateOrder updates the order of a featured cake
func (r *homeRepositoryImpl) UpdateOrder(ctx context.Context, id string, order int) (*config.FeaturedCake, error) {
	existing, err := r.GetFeaturedCakeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Order = order
	existing.UpdatedAt = time.Now()

	_, err = r.client.Collection(homeCollection).Doc("default").Collection(featuredCakesSubcollection).Doc(id).Set(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}
