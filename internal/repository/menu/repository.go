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
	GetAllCategories(ctx context.Context) ([]string, error)
	CreateCategory(ctx context.Context, categoryID string) error
	DeleteCategory(ctx context.Context, categoryID string) error
	GetAll(ctx context.Context, active *bool, category string) ([]menu.MenuItem, error)
	GetByCategory(ctx context.Context, categoryID string, active *bool) ([]menu.MenuItem, error)
	GetByID(ctx context.Context, categoryID, itemID string) (*menu.MenuItem, error)
	Create(ctx context.Context, categoryID string, item *menu.MenuItem) (*menu.MenuItem, error)
	Update(ctx context.Context, categoryID, itemID string, item *menu.MenuItem) (*menu.MenuItem, error)
	Delete(ctx context.Context, categoryID, itemID string) error
	UpdateOrder(ctx context.Context, categoryID, itemID string, order int) (*menu.MenuItem, error)
}

type menuRepositoryImpl struct {
	client *firestore.Client
}

const menuCollection = "menu"
const itemsSubcollection = "items"

func NewMenuRepository(client *firestore.Client) MenuRepository {
	return &menuRepositoryImpl{client: client}
}

func (r *menuRepositoryImpl) GetAllCategories(ctx context.Context) ([]string, error) {
	snapshots, err := r.client.Collection(menuCollection).DocumentRefs(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	var categories []string
	for _, ref := range snapshots {
		categories = append(categories, ref.ID)
	}

	return categories, nil
}

func (r *menuRepositoryImpl) CreateCategory(ctx context.Context, categoryID string) error {
	_, err := r.client.Collection(menuCollection).Doc(categoryID).Set(ctx, map[string]interface{}{
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	})
	return err
}

func (r *menuRepositoryImpl) DeleteCategory(ctx context.Context, categoryID string) error {
	items, err := r.GetByCategory(ctx, categoryID, nil)
	if err == nil {
		for _, item := range items {
			if err := r.Delete(ctx, categoryID, item.ID); err != nil {
				return err
			}
		}
	}

	_, err = r.client.Collection(menuCollection).Doc(categoryID).Delete(ctx)
	return err
}

func (r *menuRepositoryImpl) GetAll(ctx context.Context, active *bool, category string) ([]menu.MenuItem, error) {
	var items []menu.MenuItem

	if category != "" {
		return r.GetByCategory(ctx, category, active)
	}

	categories, err := r.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	for _, cat := range categories {
		catItems, err := r.GetByCategory(ctx, cat, active)
		if err == nil {
			items = append(items, catItems...)
		}
	}

	return items, nil
}

func (r *menuRepositoryImpl) GetByCategory(ctx context.Context, categoryID string, active *bool) ([]menu.MenuItem, error) {
	var items []menu.MenuItem

	docs, err := r.client.Collection(menuCollection).Doc(categoryID).Collection(itemsSubcollection).
		OrderBy("order", firestore.Asc).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch items for category %s: %w", categoryID, err)
	}

	for _, doc := range docs {
		var item menu.MenuItem
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}

		if active != nil && item.Active != *active {
			continue
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *menuRepositoryImpl) GetByID(ctx context.Context, categoryID, itemID string) (*menu.MenuItem, error) {
	doc, err := r.client.Collection(menuCollection).Doc(categoryID).Collection(itemsSubcollection).Doc(itemID).Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
				Title:      "Menu Item Not Found",
				Detail:     fmt.Sprintf("No menu item found with ID: %s in category: %s", itemID, categoryID),
				HTTPStatus: http.StatusNotFound,
				Instance:   fmt.Sprintf("https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/menu/%s/items/%s", categoryID, itemID),
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

func (r *menuRepositoryImpl) Create(ctx context.Context, categoryID string, item *menu.MenuItem) (*menu.MenuItem, error) {
	if _, err := r.client.Collection(menuCollection).Doc(categoryID).Get(ctx); err != nil {
		if status.Code(err) == codes.NotFound {
			if err := r.CreateCategory(ctx, categoryID); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	_, err := r.client.Collection(menuCollection).Doc(categoryID).Collection(itemsSubcollection).Doc(item.ID).Set(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *menuRepositoryImpl) Update(ctx context.Context, categoryID, itemID string, item *menu.MenuItem) (*menu.MenuItem, error) {
	existing, err := r.GetByID(ctx, categoryID, itemID)
	if err != nil {
		return nil, err
	}

	item.ID = itemID
	item.CreatedAt = existing.CreatedAt
	item.UpdatedAt = time.Now()

	_, err = r.client.Collection(menuCollection).Doc(categoryID).Collection(itemsSubcollection).Doc(itemID).Set(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *menuRepositoryImpl) Delete(ctx context.Context, categoryID, itemID string) error {
	_, err := r.GetByID(ctx, categoryID, itemID)
	if err != nil {
		return err
	}

	_, err = r.client.Collection(menuCollection).Doc(categoryID).Collection(itemsSubcollection).Doc(itemID).Delete(ctx)
	return err
}

func (r *menuRepositoryImpl) UpdateOrder(ctx context.Context, categoryID, itemID string, order int) (*menu.MenuItem, error) {
	existing, err := r.GetByID(ctx, categoryID, itemID)
	if err != nil {
		return nil, err
	}

	existing.Order = order
	existing.UpdatedAt = time.Now()

	_, err = r.client.Collection(menuCollection).Doc(categoryID).Collection(itemsSubcollection).Doc(itemID).Set(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// FindItemCategory searches all categories to find which one contains the item ID
func (r *menuRepositoryImpl) FindItemCategory(ctx context.Context, itemID string) (string, error) {
	categories, err := r.GetAllCategories(ctx)
	if err != nil {
		return "", err
	}

	for _, cat := range categories {
		if _, err := r.GetByID(ctx, cat, itemID); err == nil {
			return cat, nil
		}
	}

	return "", fmt.Errorf("item %s not found in any category", itemID)
}
