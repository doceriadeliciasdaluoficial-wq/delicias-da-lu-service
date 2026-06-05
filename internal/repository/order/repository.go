package order

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"delicias-da-lu-service.com/mod/internal/entity/order"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderRepository interface {
	GetAll(ctx context.Context, orderStatus string, limit int, offset int) (*order.OrderListResponse, error)
	GetByID(ctx context.Context, id string) (*order.Order, error)
	Create(ctx context.Context, ord *order.Order) (*order.Order, error)
	UpdateStatus(ctx context.Context, id string, status string) (*order.Order, error)
}

type orderRepositoryImpl struct {
	client *firestore.Client
}

func NewOrderRepository(client *firestore.Client) OrderRepository {
	return orderRepositoryImpl{
		client: client,
	}
}

func (r orderRepositoryImpl) GetAll(ctx context.Context, orderStatus string, limit int, offset int) (*order.OrderListResponse, error) {
	var orders []order.Order
	coll := r.client.Collection("orders")

	// Get all documents ordered by creation date
	docs, err := coll.OrderBy("createdAt", firestore.Desc).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	total := 0
	for _, doc := range docs {
		var ord order.Order
		if err := doc.DataTo(&ord); err != nil {
			return nil, err
		}

		// Filter by status if provided
		if orderStatus != "" && ord.Status != orderStatus {
			continue
		}

		total++

		// Apply pagination
		if total > offset && total <= (offset+limit) {
			orders = append(orders, ord)
		}
	}

	return &order.OrderListResponse{
		Total:  total,
		Limit:  limit,
		Offset: offset,
		Data:   orders,
	}, nil
}

func (r orderRepositoryImpl) GetByID(ctx context.Context, id string) (*order.Order, error) {
	doc, err := r.client.Collection("orders").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
				Title:      "Order Not Found",
				Detail:     fmt.Sprintf("No order found with ID: %s", id),
				HTTPStatus: http.StatusNotFound,
				Instance:   fmt.Sprintf("https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/orders/%s", id),
				Severity:   problemdetails.Err,
			})
		}
		return nil, err
	}

	var ord order.Order
	if err := doc.DataTo(&ord); err != nil {
		return nil, err
	}

	return &ord, nil
}

func (r orderRepositoryImpl) Create(ctx context.Context, ord *order.Order) (*order.Order, error) {
	ord.CreatedAt = time.Now()
	ord.UpdatedAt = time.Now()
	if ord.Status == "" {
		ord.Status = "pending"
	}

	documentRef := r.client.Collection("orders").NewDoc()
	ord.ID = documentRef.ID

	_, err := documentRef.Set(ctx, ord)
	if err != nil {
		return nil, err
	}

	return ord, nil
}

func (r orderRepositoryImpl) UpdateStatus(ctx context.Context, id string, newStatus string) (*order.Order, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Status = newStatus
	existing.UpdatedAt = time.Now()

	_, err = r.client.Collection("orders").Doc(id).Set(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}
