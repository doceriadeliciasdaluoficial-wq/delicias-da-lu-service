package order

import (
	"context"

	"delicias-da-lu-service.com/mod/internal/entity/order"
	orderRepo "delicias-da-lu-service.com/mod/internal/repository/order"
)

type OrderUseCase interface {
	GetAll(ctx context.Context, status string, limit int, offset int) (*order.OrderListResponse, error)
	GetByID(ctx context.Context, id string) (*order.Order, error)
	Create(ctx context.Context, ord *order.Order) (*order.Order, error)
	UpdateStatus(ctx context.Context, id string, status string) (*order.Order, error)
}

type orderUseCaseImpl struct {
	repository orderRepo.OrderRepository
}

func NewOrderUseCase(repository orderRepo.OrderRepository) OrderUseCase {
	return orderUseCaseImpl{
		repository: repository,
	}
}

func (o orderUseCaseImpl) GetAll(ctx context.Context, status string, limit int, offset int) (*order.OrderListResponse, error) {
	if limit == 0 {
		limit = 20
	}
	return o.repository.GetAll(ctx, status, limit, offset)
}

func (o orderUseCaseImpl) GetByID(ctx context.Context, id string) (*order.Order, error) {
	return o.repository.GetByID(ctx, id)
}

func (o orderUseCaseImpl) Create(ctx context.Context, ord *order.Order) (*order.Order, error) {
	return o.repository.Create(ctx, ord)
}

func (o orderUseCaseImpl) UpdateStatus(ctx context.Context, id string, status string) (*order.Order, error) {
	return o.repository.UpdateStatus(ctx, id, status)
}
