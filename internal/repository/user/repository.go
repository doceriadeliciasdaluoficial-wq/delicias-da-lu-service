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
	GetByID(ctx context.Context, id string) (*user.AdminUser, error)
	Create(ctx context.Context, usr *user.AdminUser) (*user.AdminUser, error)
	UpdateLastLogin(ctx context.Context, id string) error
}

type userRepositoryImpl struct {
	client *firestore.Client
}

func NewUserRepository(client *firestore.Client) UserRepository {
	return userRepositoryImpl{
		client: client,
	}
}

func (r userRepositoryImpl) GetByUsername(ctx context.Context, username string) (*user.AdminUser, error) {
	docs, err := r.client.Collection("users").Where("username", "==", username).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
			Title:      "User Not Found",
			Detail:     fmt.Sprintf("No user found with username: %s", username),
			HTTPStatus: http.StatusNotFound,
			Instance:   "localhost:8080/v1/auth/login",
			Severity:   problemdetails.Err,
		})
	}

	var usr user.AdminUser
	if err := docs[0].DataTo(&usr); err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r userRepositoryImpl) GetByID(ctx context.Context, id string) (*user.AdminUser, error) {
	doc, err := r.client.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
				Title:      "User Not Found",
				Detail:     fmt.Sprintf("No user found with ID: %s", id),
				HTTPStatus: http.StatusNotFound,
				Instance:   fmt.Sprintf("localhost:8080/v1/users/%s", id),
				Severity:   problemdetails.Err,
			})
		}
		return nil, err
	}

	var usr user.AdminUser
	if err := doc.DataTo(&usr); err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r userRepositoryImpl) Create(ctx context.Context, usr *user.AdminUser) (*user.AdminUser, error) {
	usr.CreatedAt = time.Now()
	usr.UpdatedAt = time.Now()

	_, err := r.client.Collection("users").Doc(usr.ID).Set(ctx, usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (r userRepositoryImpl) UpdateLastLogin(ctx context.Context, id string) error {
	_, err := r.client.Collection("users").Doc(id).Update(ctx, []firestore.Update{
		{Path: "lastLogin", Value: time.Now()},
	})
	return err
}
