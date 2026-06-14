package menu

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"delicias-da-lu-service.com/mod/internal/entity/menu"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MenuRepository interface {
	GetAll(ctx context.Context, active *bool, category string) ([]menu.MenuItem, error)
	GetByID(ctx context.Context, id string) (*menu.MenuItem, error)
	Create(ctx context.Context, item *menu.MenuItem) (*menu.MenuItem, error)
	Update(ctx context.Context, id string, item *menu.MenuItem) (*menu.MenuItem, error)
	Delete(ctx context.Context, id string) error
	UpdateOrder(ctx context.Context, id string, order int) (*menu.MenuItem, error)
}

type menuRepositoryImpl struct {
	client *firestore.Client
}

func NewMenuRepository(client *firestore.Client) MenuRepository {
	return menuRepositoryImpl{
		client: client,
	}
}

func (r menuRepositoryImpl) GetAll(ctx context.Context, active *bool, category string) ([]menu.MenuItem, error) {
	var items []menu.MenuItem
	coll := r.client.Collection("menu")

	docs, err := coll.OrderBy("order", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		var item menu.MenuItem
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}

		skip := false
		if active != nil && item.Active != *active {
			skip = true
		}
		if category != "" && item.Category != category {
			skip = true
		}

		if !skip {
			items = append(items, item)
		}
	}

	return items, nil
}

func (r menuRepositoryImpl) GetByID(ctx context.Context, id string) (*menu.MenuItem, error) {
	doc, err := r.client.Collection("menu").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
				Title:      "Menu Item Not Found",
				Detail:     fmt.Sprintf("No menu item found with ID: %s", id),
				HTTPStatus: http.StatusNotFound,
				Instance:   fmt.Sprintf("https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/menu/items/%s", id),
				Severity:   problemdetails.Err,
			})
		}
		return nil, err
	}

	var item menu.MenuItem
	if err := doc.DataTo(&item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r menuRepositoryImpl) Create(ctx context.Context, item *menu.MenuItem) (*menu.MenuItem, error) {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	_, err := r.client.Collection("menu").Doc(item.ID).Set(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r menuRepositoryImpl) Update(ctx context.Context, id string, item *menu.MenuItem) (*menu.MenuItem, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	item.ID = id
	item.CreatedAt = existing.CreatedAt
	item.UpdatedAt = time.Now()

	_, err = r.client.Collection("menu").Doc(id).Set(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r menuRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = r.client.Collection("menu").Doc(id).Delete(ctx)
	return err
}

func (r menuRepositoryImpl) UpdateOrder(ctx context.Context, id string, order int) (*menu.MenuItem, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Order = order
	existing.UpdatedAt = time.Now()

	_, err = r.client.Collection("menu").Doc(id).Set(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}
