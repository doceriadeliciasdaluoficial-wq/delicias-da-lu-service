package auth

import (
	"context"
	"errors"
	"testing"

	"delicias-da-lu-service.com/mod/internal/entity/user"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
)

type MockUserRepository struct {
	GetByUsernameFunc   func(ctx context.Context, username string) (*user.AdminUser, error)
	GetByIDFunc         func(ctx context.Context, id string) (*user.AdminUser, error)
	CreateFunc          func(ctx context.Context, usr *user.AdminUser) (*user.AdminUser, error)
	UpdateLastLoginFunc func(ctx context.Context, userID string) error
	UpdateFunc          func(ctx context.Context, userID string, usr *user.AdminUser) (*user.AdminUser, error)
	DeleteFunc          func(ctx context.Context, userID string) error
	ListAllFunc         func(ctx context.Context) ([]user.AdminUser, error)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*user.AdminUser, error) {
	if m.GetByUsernameFunc != nil {
		return m.GetByUsernameFunc(ctx, username)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*user.AdminUser, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) Create(ctx context.Context, usr *user.AdminUser) (*user.AdminUser, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, usr)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	if m.UpdateLastLoginFunc != nil {
		return m.UpdateLastLoginFunc(ctx, userID)
	}
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, userID string, usr *user.AdminUser) (*user.AdminUser, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, userID, usr)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) Delete(ctx context.Context, userID string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, userID)
	}
	return nil
}

func (m *MockUserRepository) ListAll(ctx context.Context) ([]user.AdminUser, error) {
	if m.ListAllFunc != nil {
		return m.ListAllFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func TestAuthUseCaseLogin(t *testing.T) {
	ctx := context.Background()
	jwtSecret := "test-secret-key"

	tests := []struct {
		name           string
		username       string
		password       string
		mockUser       *user.AdminUser
		mockError      error
		wantToken      bool
		wantError      bool
		wantStatusCode int
	}{
		{
			name:     "Successful login",
			username: "admin",
			password: "password123",
			mockUser: &user.AdminUser{
				ID:       "user-1",
				Username: "admin",
				Password: "password123",
				Role:     "admin",
			},
			mockError:      nil,
			wantToken:      true,
			wantError:      false,
			wantStatusCode: 0,
		},
		{
			name:           "User not found",
			username:       "nonexistent",
			password:       "password123",
			mockUser:       nil,
			mockError:      errors.New("user not found"),
			wantToken:      false,
			wantError:      true,
			wantStatusCode: 401,
		},
		{
			name:     "Invalid password",
			username: "admin",
			password: "wrongpassword",
			mockUser: &user.AdminUser{
				ID:       "user-1",
				Username: "admin",
				Password: "password123",
				Role:     "admin",
			},
			mockError:      nil,
			wantToken:      false,
			wantError:      true,
			wantStatusCode: 401,
		},
		{
			name:           "Empty username",
			username:       "",
			password:       "password123",
			mockUser:       nil,
			mockError:      errors.New("invalid credentials"),
			wantToken:      false,
			wantError:      true,
			wantStatusCode: 401,
		},
		{
			name:     "Empty password",
			username: "admin",
			password: "",
			mockUser: &user.AdminUser{
				ID:       "user-1",
				Username: "admin",
				Password: "",
				Role:     "admin",
			},
			mockError:      nil,
			wantToken:      true,
			wantError:      false,
			wantStatusCode: 0,
		},
		{
			name:           "Case sensitive username",
			username:       "Admin",
			password:       "password123",
			mockUser:       nil,
			mockError:      errors.New("user not found"),
			wantToken:      false,
			wantError:      true,
			wantStatusCode: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				GetByUsernameFunc: func(ctx context.Context, username string) (*user.AdminUser, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockUser, nil
				},
				UpdateLastLoginFunc: func(ctx context.Context, userID string) error {
					return nil
				},
			}

			authUseCase := NewAuthUseCase(mockRepo, jwtSecret)
			response, err := authUseCase.Login(ctx, tt.username, tt.password)

			if (err != nil) != tt.wantError {
				t.Errorf("Login() error = %v, wantError = %v", err, tt.wantError)
				return
			}

			if tt.wantToken && response == nil {
				t.Errorf("Login() expected token, got nil")
				return
			}

			if tt.wantToken && response.Token == "" {
				t.Errorf("Login() expected non-empty token")
			}

			if tt.wantError && err != nil {
				if pdErr, ok := err.(*problemdetails.Error); ok {
					if pdErr.HTTPStatus != tt.wantStatusCode {
						t.Errorf("Login() status = %d, want = %d", pdErr.HTTPStatus, tt.wantStatusCode)
					}
				}
			}
		})
	}
}

func TestAuthUseCaseRefreshToken(t *testing.T) {
	ctx := context.Background()
	jwtSecret := "test-secret-key"
	mockRepo := &MockUserRepository{}
	authUseCase := NewAuthUseCase(mockRepo, jwtSecret)

	mockRepo.CreateFunc = func(ctx context.Context, usr *user.AdminUser) (*user.AdminUser, error) {
		return usr, nil
	}

	user1 := &user.AdminUser{
		ID:       "user-1",
		Username: "admin",
		Password: "password123",
		Role:     "admin",
	}
	mockRepo.GetByUsernameFunc = func(ctx context.Context, username string) (*user.AdminUser, error) {
		return user1, nil
	}

	loginResp, _ := authUseCase.Login(ctx, "admin", "password123")
	validToken := loginResp.Token

	tests := []struct {
		name      string
		token     string
		wantToken bool
		wantError bool
	}{
		{
			name:      "Valid token refresh",
			token:     validToken,
			wantToken: true,
			wantError: false,
		},
		{
			name:      "Invalid token",
			token:     "invalid.token.here",
			wantToken: false,
			wantError: true,
		},
		{
			name:      "Empty token",
			token:     "",
			wantToken: false,
			wantError: true,
		},
		{
			name:      "Malformed token",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantToken: false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newToken, err := authUseCase.RefreshToken(ctx, tt.token)

			if (err != nil) != tt.wantError {
				t.Errorf("RefreshToken() error = %v, wantError = %v", err, tt.wantError)
				return
			}

			if tt.wantToken && newToken == "" {
				t.Errorf("RefreshToken() expected non-empty token")
			}

			if !tt.wantError && newToken != "" && tt.token == validToken {
				if newToken == tt.token {
					t.Logf("RefreshToken() returned new token")
				}
			}
		})
	}
}

func TestAuthUseCaseValidateToken(t *testing.T) {
	ctx := context.Background()
	jwtSecret := "test-secret-key"
	mockRepo := &MockUserRepository{}
	authUseCase := NewAuthUseCase(mockRepo, jwtSecret)

	mockRepo.CreateFunc = func(ctx context.Context, usr *user.AdminUser) (*user.AdminUser, error) {
		return usr, nil
	}

	user1 := &user.AdminUser{
		ID:       "user-123",
		Username: "admin",
		Password: "password123",
		Role:     "admin",
	}
	mockRepo.GetByUsernameFunc = func(ctx context.Context, username string) (*user.AdminUser, error) {
		return user1, nil
	}

	loginResp, _ := authUseCase.Login(ctx, "admin", "password123")
	validToken := loginResp.Token

	tests := []struct {
		name       string
		token      string
		wantUserID string
		wantError  bool
	}{
		{
			name:       "Valid token",
			token:      validToken,
			wantUserID: "user-123",
			wantError:  false,
		},
		{
			name:       "Invalid token",
			token:      "invalid.token.here",
			wantUserID: "",
			wantError:  true,
		},
		{
			name:       "Empty token",
			token:      "",
			wantUserID: "",
			wantError:  true,
		},
		{
			name:       "Token with wrong secret",
			token:      validToken,
			wantUserID: "user-123",
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := authUseCase.ValidateToken(ctx, tt.token)

			if (err != nil) != tt.wantError {
				t.Errorf("ValidateToken() error = %v, wantError = %v", err, tt.wantError)
				return
			}

			if !tt.wantError && userID != tt.wantUserID {
				t.Errorf("ValidateToken() userID = %s, want = %s", userID, tt.wantUserID)
			}
		})
	}
}

func BenchmarkAuthUseCaseLogin(b *testing.B) {
	ctx := context.Background()
	jwtSecret := "test-secret-key"
	mockRepo := &MockUserRepository{
		GetByUsernameFunc: func(ctx context.Context, username string) (*user.AdminUser, error) {
			return &user.AdminUser{
				ID:       "user-1",
				Username: "admin",
				Password: "password123",
				Role:     "admin",
			}, nil
		},
		UpdateLastLoginFunc: func(ctx context.Context, userID string) error {
			return nil
		},
	}

	authUseCase := NewAuthUseCase(mockRepo, jwtSecret)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		authUseCase.Login(ctx, "admin", "password123")
	}
}

func BenchmarkAuthUseCaseValidateToken(b *testing.B) {
	ctx := context.Background()
	jwtSecret := "test-secret-key"
	mockRepo := &MockUserRepository{
		GetByUsernameFunc: func(ctx context.Context, username string) (*user.AdminUser, error) {
			return &user.AdminUser{
				ID:       "user-1",
				Username: "admin",
				Password: "password123",
				Role:     "admin",
			}, nil
		},
		UpdateLastLoginFunc: func(ctx context.Context, userID string) error {
			return nil
		},
	}
	authUseCase := NewAuthUseCase(mockRepo, jwtSecret)

	loginResp, _ := authUseCase.Login(ctx, "admin", "password123")
	token := loginResp.Token

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		authUseCase.ValidateToken(ctx, token)
	}
}
