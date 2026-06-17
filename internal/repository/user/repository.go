package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"delicias-da-lu-service.com/mod/internal/entity/user"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*user.AdminUser, error)
	GetByID(ctx context.Context, userID string) (*user.AdminUser, error)
	Create(ctx context.Context, adminUser *user.AdminUser) (*user.AdminUser, error)
	Update(ctx context.Context, userID string, adminUser *user.AdminUser) (*user.AdminUser, error)
	UpdateLastLogin(ctx context.Context, userID string) error
	Delete(ctx context.Context, userID string) error
	ListAll(ctx context.Context) ([]user.AdminUser, error)
}

type userRepositoryImpl struct {
	client *firestore.Client
}

const adminCollection = "admin"
const usersSubcollection = "users"

func NewUserRepository(client *firestore.Client) UserRepository {
	return &userRepositoryImpl{client: client}
}

func (r *userRepositoryImpl) GetByUsername(ctx context.Context, username string) (*user.AdminUser, error) {
	docs, err := r.client.Collection(adminCollection).Doc("default").Collection(usersSubcollection).
		Where("username", "==", username).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if len(docs) == 0 {
		return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
			Title:      "User Not Found",
			Detail:     fmt.Sprintf("No user found with username: %s", username),
			HTTPStatus: http.StatusNotFound,
			Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/auth/login",
			Severity:   problemdetails.Err,
		})
	}

	var adminUser user.AdminUser
	if err := docs[0].DataTo(&adminUser); err != nil {
		return nil, err
	}

	return &adminUser, nil
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, userID string) (*user.AdminUser, error) {
	doc, err := r.client.Collection(adminCollection).Doc("default").Collection(usersSubcollection).Doc(userID).Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
				Title:      "User Not Found",
				Detail:     fmt.Sprintf("No user found with ID: %s", userID),
				HTTPStatus: http.StatusNotFound,
				Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/auth",
				Severity:   problemdetails.Err,
			})
		}
		return nil, err
	}

	var adminUser user.AdminUser
	if err := doc.DataTo(&adminUser); err != nil {
		return nil, err
	}

	return &adminUser, nil
}

func (r *userRepositoryImpl) Create(ctx context.Context, adminUser *user.AdminUser) (*user.AdminUser, error) {
	if _, err := r.client.Collection(adminCollection).Doc("default").Get(ctx); err != nil {
		if status.Code(err) == codes.NotFound {
			if _, err := r.client.Collection(adminCollection).Doc("default").Set(ctx, map[string]interface{}{
				"createdAt": time.Now(),
				"updatedAt": time.Now(),
			}); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	adminUser.CreatedAt = time.Now()
	adminUser.UpdatedAt = time.Now()
	now := time.Now()
	adminUser.LastLogin = now

	_, err := r.client.Collection(adminCollection).Doc("default").Collection(usersSubcollection).Doc(adminUser.ID).Set(ctx, adminUser)
	if err != nil {
		return nil, err
	}

	return adminUser, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, userID string, adminUser *user.AdminUser) (*user.AdminUser, error) {
	existing, err := r.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	adminUser.ID = userID
	adminUser.CreatedAt = existing.CreatedAt
	adminUser.UpdatedAt = time.Now()

	_, err = r.client.Collection(adminCollection).Doc("default").Collection(usersSubcollection).Doc(userID).Set(ctx, adminUser)
	if err != nil {
		return nil, err
	}

	return adminUser, nil
}

func (r *userRepositoryImpl) UpdateLastLogin(ctx context.Context, userID string) error {
	existing, err := r.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now()
	existing.LastLogin = now
	existing.UpdatedAt = time.Now()

	_, err = r.client.Collection(adminCollection).Doc("default").Collection(usersSubcollection).Doc(userID).Set(ctx, existing)
	return err
}

func (r *userRepositoryImpl) Delete(ctx context.Context, userID string) error {
	_, err := r.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = r.client.Collection(adminCollection).Doc("default").Collection(usersSubcollection).Doc(userID).Delete(ctx)
	return err
}

func (r *userRepositoryImpl) ListAll(ctx context.Context) ([]user.AdminUser, error) {
	var users []user.AdminUser

	docs, err := r.client.Collection(adminCollection).Doc("default").Collection(usersSubcollection).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	for _, doc := range docs {
		var adminUser user.AdminUser
		if err := doc.DataTo(&adminUser); err != nil {
			return nil, err
		}
		users = append(users, adminUser)
	}

	return users, nil
}
