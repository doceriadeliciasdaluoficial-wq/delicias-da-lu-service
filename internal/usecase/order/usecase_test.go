package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"delicias-da-lu-service.com/mod/internal/entity/order"
)

type MockOrderRepository struct {
	GetAllFunc       func(ctx context.Context, status string, limit int, offset int) (*order.OrderListResponse, error)
	GetByIDFunc      func(ctx context.Context, id string) (*order.Order, error)
	CreateFunc       func(ctx context.Context, ord *order.Order) (*order.Order, error)
	UpdateStatusFunc func(ctx context.Context, id string, status string) (*order.Order, error)
}

func (m *MockOrderRepository) GetAll(ctx context.Context, status string, limit int, offset int) (*order.OrderListResponse, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, status, limit, offset)
	}
	return nil, errors.New("not implemented")
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id string) (*order.Order, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockOrderRepository) Create(ctx context.Context, ord *order.Order) (*order.Order, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, ord)
	}
	return nil, errors.New("not implemented")
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id string, status string) (*order.Order, error) {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil, errors.New("not implemented")
}

func TestOrderUseCaseGetAll(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		status     string
		limit      int
		offset     int
		mockOrders []order.Order
		mockTotal  int
		mockError  error
		wantTotal  int
		wantLimit  int
		wantOffset int
		wantError  bool
	}{
		{
			name:   "Get all orders with default limit",
			status: "",
			limit:  0,
			offset: 0,
			mockOrders: []order.Order{
				{ID: "order-1", Status: "pending"},
			},
			mockTotal:  1,
			mockError:  nil,
			wantTotal:  1,
			wantLimit:  20,
			wantOffset: 0,
			wantError:  false,
		},
		{
			name:   "Get orders with custom limit",
			status: "",
			limit:  10,
			offset: 0,
			mockOrders: []order.Order{
				{ID: "order-1", Status: "pending"},
			},
			mockTotal:  1,
			mockError:  nil,
			wantTotal:  1,
			wantLimit:  10,
			wantOffset: 0,
			wantError:  false,
		},
		{
			name:   "Get orders with offset",
			status: "",
			limit:  10,
			offset: 20,
			mockOrders: []order.Order{
				{ID: "order-21", Status: "completed"},
			},
			mockTotal:  100,
			mockError:  nil,
			wantTotal:  100,
			wantLimit:  10,
			wantOffset: 20,
			wantError:  false,
		},
		{
			name:   "Get orders filtered by status",
			status: "pending",
			limit:  0,
			offset: 0,
			mockOrders: []order.Order{
				{ID: "order-1", Status: "pending"},
			},
			mockTotal:  1,
			mockError:  nil,
			wantTotal:  1,
			wantLimit:  20,
			wantOffset: 0,
			wantError:  false,
		},
		{
			name:       "Get orders repository error",
			status:     "",
			limit:      10,
			offset:     0,
			mockOrders: nil,
			mockTotal:  0,
			mockError:  errors.New("database error"),
			wantError:  true,
		},
		{
			name:       "Empty result",
			status:     "",
			limit:      10,
			offset:     0,
			mockOrders: []order.Order{},
			mockTotal:  0,
			mockError:  nil,
			wantTotal:  0,
			wantLimit:  10,
			wantOffset: 0,
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockOrderRepository{
				GetAllFunc: func(ctx context.Context, status string, limit int, offset int) (*order.OrderListResponse, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return &order.OrderListResponse{
						Total:  tt.mockTotal,
						Limit:  limit,
						Offset: offset,
						Data:   tt.mockOrders,
					}, nil
				},
			}

			useCase := NewOrderUseCase(mockRepo)
			response, err := useCase.GetAll(ctx, tt.status, tt.limit, tt.offset)

			if (err != nil) != tt.wantError {
				t.Errorf("GetAll() error = %v, wantError = %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if response.Total != tt.wantTotal {
					t.Errorf("GetAll() Total = %d, want = %d", response.Total, tt.wantTotal)
				}
				if response.Limit != tt.wantLimit {
					t.Errorf("GetAll() Limit = %d, want = %d", response.Limit, tt.wantLimit)
				}
				if response.Offset != tt.wantOffset {
					t.Errorf("GetAll() Offset = %d, want = %d", response.Offset, tt.wantOffset)
				}
			}
		})
	}
}

func TestOrderUseCaseGetByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		orderID   string
		mockOrder *order.Order
		mockError error
		wantError bool
	}{
		{
			name:    "Get existing order",
			orderID: "order-1",
			mockOrder: &order.Order{
				ID:         "order-1",
				Status:     "pending",
				TotalPrice: 100.00,
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Get non-existent order",
			orderID:   "non-existent",
			mockOrder: nil,
			mockError: errors.New("order not found"),
			wantError: true,
		},
		{
			name:      "Empty order ID",
			orderID:   "",
			mockOrder: nil,
			mockError: errors.New("invalid order ID"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockOrderRepository{
				GetByIDFunc: func(ctx context.Context, id string) (*order.Order, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockOrder, nil
				},
			}

			useCase := NewOrderUseCase(mockRepo)
			result, err := useCase.GetByID(ctx, tt.orderID)

			if (err != nil) != tt.wantError {
				t.Errorf("GetByID() error = %v, wantError = %v", err, tt.wantError)
				return
			}

			if !tt.wantError && result.ID != tt.mockOrder.ID {
				t.Errorf("GetByID() ID = %s, want = %s", result.ID, tt.mockOrder.ID)
			}
		})
	}
}

func TestOrderUseCaseCreate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name      string
		order     *order.Order
		mockError error
		wantError bool
	}{
		{
			name: "Create valid order",
			order: &order.Order{
				ID: "order-new",
				Items: []order.OrderItem{
					{
						Type:      "menu",
						Quantity:  2,
						UnitPrice: 50.00,
						Subtotal:  100.00,
					},
				},
				CustomerInfo: order.CustomerInfo{
					Name:  "John",
					Phone: "123456789",
				},
				Status:     "pending",
				TotalPrice: 100.00,
				CreatedAt:  now,
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "Create order with empty items",
			order: &order.Order{
				ID:    "order-empty",
				Items: []order.OrderItem{},
				CustomerInfo: order.CustomerInfo{
					Name:  "Jane",
					Phone: "987654321",
				},
				Status:     "pending",
				TotalPrice: 0,
				CreatedAt:  now,
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Create with repository error",
			order:     &order.Order{ID: "order-fail"},
			mockError: errors.New("database error"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockOrderRepository{
				CreateFunc: func(ctx context.Context, ord *order.Order) (*order.Order, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return ord, nil
				},
			}

			useCase := NewOrderUseCase(mockRepo)
			result, err := useCase.Create(ctx, tt.order)

			if (err != nil) != tt.wantError {
				t.Errorf("Create() error = %v, wantError = %v", err, tt.wantError)
				return
			}

			if !tt.wantError && result.ID != tt.order.ID {
				t.Errorf("Create() ID = %s, want = %s", result.ID, tt.order.ID)
			}
		})
	}
}

func TestOrderUseCaseUpdateStatus(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		orderID   string
		newStatus string
		mockOrder *order.Order
		mockError error
		wantError bool
	}{
		{
			name:      "Update to confirmed",
			orderID:   "order-1",
			newStatus: "confirmed",
			mockOrder: &order.Order{
				ID:     "order-1",
				Status: "confirmed",
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Update to preparing",
			orderID:   "order-1",
			newStatus: "preparing",
			mockOrder: &order.Order{
				ID:     "order-1",
				Status: "preparing",
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Update to delivered",
			orderID:   "order-1",
			newStatus: "delivered",
			mockOrder: &order.Order{
				ID:     "order-1",
				Status: "delivered",
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Update to cancelled",
			orderID:   "order-1",
			newStatus: "cancelled",
			mockOrder: &order.Order{
				ID:     "order-1",
				Status: "cancelled",
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Update non-existent order",
			orderID:   "non-existent",
			newStatus: "confirmed",
			mockOrder: nil,
			mockError: errors.New("order not found"),
			wantError: true,
		},
		{
			name:      "Update with invalid status",
			orderID:   "order-1",
			newStatus: "invalid_status",
			mockOrder: &order.Order{
				ID:     "order-1",
				Status: "invalid_status",
			},
			mockError: nil,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockOrderRepository{
				UpdateStatusFunc: func(ctx context.Context, id string, status string) (*order.Order, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockOrder, nil
				},
			}

			useCase := NewOrderUseCase(mockRepo)
			result, err := useCase.UpdateStatus(ctx, tt.orderID, tt.newStatus)

			if (err != nil) != tt.wantError {
				t.Errorf("UpdateStatus() error = %v, wantError = %v", err, tt.wantError)
				return
			}

			if !tt.wantError && result.Status != tt.newStatus {
				t.Errorf("UpdateStatus() Status = %s, want = %s", result.Status, tt.newStatus)
			}
		})
	}
}

func BenchmarkOrderUseCaseGetAll(b *testing.B) {
	ctx := context.Background()
	mockRepo := &MockOrderRepository{
		GetAllFunc: func(ctx context.Context, status string, limit int, offset int) (*order.OrderListResponse, error) {
			return &order.OrderListResponse{
				Total:  100,
				Limit:  limit,
				Offset: offset,
				Data: []order.Order{
					{ID: "order-1", Status: "pending"},
				},
			}, nil
		},
	}

	useCase := NewOrderUseCase(mockRepo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		useCase.GetAll(ctx, "", 20, 0)
	}
}

func BenchmarkOrderUseCaseGetByID(b *testing.B) {
	ctx := context.Background()
	mockRepo := &MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*order.Order, error) {
			return &order.Order{
				ID:     "order-1",
				Status: "pending",
			}, nil
		},
	}

	useCase := NewOrderUseCase(mockRepo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		useCase.GetByID(ctx, "order-1")
	}
}
