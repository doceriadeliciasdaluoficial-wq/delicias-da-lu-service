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
	GetAllTypes(ctx context.Context) ([]string, error)
	CreateType(ctx context.Context, typeID string) error
	DeleteType(ctx context.Context, typeID string) error
	GetByType(ctx context.Context, componentType string, active *bool) ([]cakebuilder.CakeBuilderComponent, error)
	GetByID(ctx context.Context, componentType, id string) (*cakebuilder.CakeBuilderComponent, error)
	Create(ctx context.Context, componentType string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error)
	Update(ctx context.Context, componentType, id string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error)
	Delete(ctx context.Context, componentType, id string) error
	GetAll(ctx context.Context, active *bool) (map[string][]cakebuilder.CakeBuilderComponent, error)
	UpdateOrder(ctx context.Context, componentType, id string, order int) (*cakebuilder.CakeBuilderComponent, error)
}

type cakeBuilderRepositoryImpl struct {
	client *firestore.Client
}

const cakeBuilderCollection = "cakeBuilder"
const componentsSubcollection = "components"

var validTypes = []string{"massas", "recheios", "coberturas", "decoracoes", "sizes"}

func NewCakeBuilderRepository(client *firestore.Client) CakeBuilderRepository {
	return &cakeBuilderRepositoryImpl{client: client}
}

func (r *cakeBuilderRepositoryImpl) GetAllTypes(ctx context.Context) ([]string, error) {
	snapshots, err := r.client.Collection(cakeBuilderCollection).DocumentRefs(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch types: %w", err)
	}

	var types []string
	for _, ref := range snapshots {
		types = append(types, ref.ID)
	}

	return types, nil
}

func (r *cakeBuilderRepositoryImpl) CreateType(ctx context.Context, typeID string) error {
	_, err := r.client.Collection(cakeBuilderCollection).Doc(typeID).Set(ctx, map[string]interface{}{
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	})
	return err
}

func (r *cakeBuilderRepositoryImpl) DeleteType(ctx context.Context, typeID string) error {
	components, err := r.GetByType(ctx, typeID, nil)
	if err == nil {
		for _, comp := range components {
			if err := r.Delete(ctx, typeID, comp.ID); err != nil {
				return err
			}
		}
	}

	_, err = r.client.Collection(cakeBuilderCollection).Doc(typeID).Delete(ctx)
	return err
}

func (r *cakeBuilderRepositoryImpl) GetByType(ctx context.Context, componentType string, active *bool) ([]cakebuilder.CakeBuilderComponent, error) {
	var components []cakebuilder.CakeBuilderComponent

	docs, err := r.client.Collection(cakeBuilderCollection).Doc(componentType).Collection(componentsSubcollection).
		OrderBy("order", firestore.Asc).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch components for type %s: %w", componentType, err)
	}

	for _, doc := range docs {
		var component cakebuilder.CakeBuilderComponent
		if err := doc.DataTo(&component); err != nil {
			return nil, err
		}

		if active != nil && component.Active != *active {
			continue
		}

		component.Type = componentType
		components = append(components, component)
	}

	return components, nil
}

func (r *cakeBuilderRepositoryImpl) GetByID(ctx context.Context, componentType, id string) (*cakebuilder.CakeBuilderComponent, error) {
	doc, err := r.client.Collection(cakeBuilderCollection).Doc(componentType).Collection(componentsSubcollection).Doc(id).Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
				Title:      "Cake Builder Component Not Found",
				Detail:     fmt.Sprintf("No component found with ID: %s of type: %s", id, componentType),
				HTTPStatus: http.StatusNotFound,
				Instance:   fmt.Sprintf("https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/cake-builder/%s/%s", componentType, id),
				Severity:   problemdetails.Err,
			})
		}
		return nil, err
	}

	var component cakebuilder.CakeBuilderComponent
	if err := doc.DataTo(&component); err != nil {
		return nil, err
	}

	component.Type = componentType
	return &component, nil
}

func (r *cakeBuilderRepositoryImpl) Create(ctx context.Context, componentType string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error) {
	if _, err := r.client.Collection(cakeBuilderCollection).Doc(componentType).Get(ctx); err != nil {
		if status.Code(err) == codes.NotFound {
			if err := r.CreateType(ctx, componentType); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	component.Type = componentType
	component.CreatedAt = time.Now()
	component.UpdatedAt = time.Now()

	_, err := r.client.Collection(cakeBuilderCollection).Doc(componentType).Collection(componentsSubcollection).Doc(component.ID).Set(ctx, component)
	if err != nil {
		return nil, err
	}

	return component, nil
}

func (r *cakeBuilderRepositoryImpl) Update(ctx context.Context, componentType, id string, component *cakebuilder.CakeBuilderComponent) (*cakebuilder.CakeBuilderComponent, error) {
	existing, err := r.GetByID(ctx, componentType, id)
	if err != nil {
		return nil, err
	}

	component.ID = id
	component.Type = componentType
	component.CreatedAt = existing.CreatedAt
	component.UpdatedAt = time.Now()

	_, err = r.client.Collection(cakeBuilderCollection).Doc(componentType).Collection(componentsSubcollection).Doc(id).Set(ctx, component)
	if err != nil {
		return nil, err
	}

	return component, nil
}

func (r *cakeBuilderRepositoryImpl) Delete(ctx context.Context, componentType, id string) error {
	_, err := r.GetByID(ctx, componentType, id)
	if err != nil {
		return err
	}

	_, err = r.client.Collection(cakeBuilderCollection).Doc(componentType).Collection(componentsSubcollection).Doc(id).Delete(ctx)
	return err
}

func (r *cakeBuilderRepositoryImpl) GetAll(ctx context.Context, active *bool) (map[string][]cakebuilder.CakeBuilderComponent, error) {
	result := make(map[string][]cakebuilder.CakeBuilderComponent)

	for _, t := range validTypes {
		components, err := r.GetByType(ctx, t, active)
		if err != nil {
			return nil, err
		}
		if components == nil {
			components = []cakebuilder.CakeBuilderComponent{}
		}
		result[t] = components
	}

	return result, nil
}

func (r *cakeBuilderRepositoryImpl) UpdateOrder(ctx context.Context, componentType, id string, order int) (*cakebuilder.CakeBuilderComponent, error) {
	existing, err := r.GetByID(ctx, componentType, id)
	if err != nil {
		return nil, err
	}

	existing.Order = order
	existing.UpdatedAt = time.Now()

	_, err = r.client.Collection(cakeBuilderCollection).Doc(componentType).Collection(componentsSubcollection).Doc(id).Set(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}
