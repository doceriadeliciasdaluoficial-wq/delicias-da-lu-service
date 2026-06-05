package cakebuilder

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"delicias-da-lu-service.com/mod/internal/entity/cakebuilder"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CakeBuilderRepository interface {
	GetByType(ctx context.Context, componentType string, active *bool) ([]cakebuilder.CakeBuilderComponent, error)
	GetByID(ctx context.Context, componentType, id string) (*cakebuilder.CakeBuilderComponent, error)
	Create(ctx context.Context, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error)
	Update(ctx context.Context, id string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error)
	Delete(ctx context.Context, componentType, id string) error
	GetAll(ctx context.Context, active *bool) (map[string][]cakebuilder.CakeBuilderComponent, error)
}

type cakeBuilderRepositoryImpl struct {
	client *firestore.Client
}

func NewCakeBuilderRepository(client *firestore.Client) CakeBuilderRepository {
	return cakeBuilderRepositoryImpl{
		client: client,
	}
}

func (r cakeBuilderRepositoryImpl) GetByType(ctx context.Context, componentType string, active *bool) ([]cakebuilder.CakeBuilderComponent, error) {
	var components []cakebuilder.CakeBuilderComponent
	coll := r.client.Collection("cakeBuilder")

	// Get all documents and filter in memory
	docs, err := coll.OrderBy("order", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		var component cakebuilder.CakeBuilderComponent
		if err := doc.DataTo(&component); err != nil {
			return nil, err
		}

		// Apply filters
		if component.Type != componentType {
			continue
		}
		if active != nil && component.Active != *active {
			continue
		}

		components = append(components, component)
	}

	return components, nil
}

func (r cakeBuilderRepositoryImpl) GetByID(ctx context.Context, componentType, id string) (*cakebuilder.CakeBuilderComponent, error) {
	doc, err := r.client.Collection("cakeBuilder").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
				Title:      "Cake Builder Component Not Found",
				Detail:     fmt.Sprintf("No component found with ID: %s", id),
				HTTPStatus: http.StatusNotFound,
				Instance:   fmt.Sprintf("localhost:8080/v1/cake-builder/%s/%s", componentType, id),
				Severity:   problemdetails.Err,
			})
		}
		return nil, err
	}

	var component cakebuilder.CakeBuilderComponent
	if err := doc.DataTo(&component); err != nil {
		return nil, err
	}

	if component.Type != componentType {
		return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/invalid-type",
			Title:      "Invalid Component Type",
			Detail:     fmt.Sprintf("Component type mismatch: expected %s, got %s", componentType, component.Type),
			HTTPStatus: http.StatusBadRequest,
			Instance:   fmt.Sprintf("localhost:8080/v1/cake-builder/%s/%s", componentType, id),
			Severity:   problemdetails.Err,
		})
	}

	return &component, nil
}

func (r cakeBuilderRepositoryImpl) Create(ctx context.Context, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error) {
	component.CreatedAt = time.Now()
	component.UpdatedAt = time.Now()

	_, err := r.client.Collection("cakeBuilder").Doc(component.ID).Set(ctx, component)
	if err != nil {
		return nil, err
	}

	return component, nil
}

func (r cakeBuilderRepositoryImpl) Update(ctx context.Context, id string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error) {
	existing, err := r.GetByID(ctx, component.Type, id)
	if err != nil {
		return nil, err
	}

	component.ID = id
	component.CreatedAt = existing.CreatedAt
	component.UpdatedAt = time.Now()

	_, err = r.client.Collection("cakeBuilder").Doc(id).Set(ctx, component)
	if err != nil {
		return nil, err
	}

	return component, nil
}

func (r cakeBuilderRepositoryImpl) Delete(ctx context.Context, componentType, id string) error {
	_, err := r.GetByID(ctx, componentType, id)
	if err != nil {
		return err
	}

	_, err = r.client.Collection("cakeBuilder").Doc(id).Delete(ctx)
	return err
}

func (r cakeBuilderRepositoryImpl) GetAll(ctx context.Context, active *bool) (map[string][]cakebuilder.CakeBuilderComponent, error) {
	types := []string{"massa", "recheio", "cobertura", "decoracao"}
	result := make(map[string][]cakebuilder.CakeBuilderComponent)

	for _, t := range types {
		components, err := r.GetByType(ctx, t, active)
		if err != nil {
			return nil, err
		}
		if components == nil {
			components = []cakebuilder.CakeBuilderComponent{}
		}
		result[t+"s"] = components
	}

	return result, nil
}
