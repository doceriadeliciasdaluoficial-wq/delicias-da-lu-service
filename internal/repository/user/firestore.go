package user

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"delicias-da-lu-service.com/mod/internal/entity/user"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepository interface {
	Create(context.Context, *user.User) (*user.User, error)
	Get(context.Context, string, string) ([]user.User, error)
	Update(context.Context, string, *user.User) (*user.User, error)
	Delete(context.Context, string) error
}

type userRepositoryImple struct {
	client *firestore.Client
}

func NewUserRepository(client *firestore.Client) UserRepository {
	return userRepositoryImple{
		client: client,
	}
}

func (ur userRepositoryImple) Create(ctx context.Context, u *user.User) (*user.User, error) {
	documentRef, _, err := ur.client.Collection("users").Add(ctx, u)
	if err != nil {
		return nil, err
	}
	u.Id = documentRef.ID
	return u, nil
}

func (ur userRepositoryImple) Get(ctx context.Context, field string, value string) ([]user.User, error) {
	var users []user.User

	iter := ur.client.Collection("users").Where(field, "==", value).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var u user.User
		if err := doc.DataTo(&u); err != nil {
			return nil, err
		}
		u.Id = doc.Ref.ID
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/user-not-found-on-get",
			Title:      "Error Type Not Found On Get",
			Detail:     "No user found for the provided field and value when searching to get",
			HTTPStatus: http.StatusNotFound,
			Instance:   "localhost:8080/v1/error/type=instance&identifier=",
			Severity:   problemdetails.Err,
		})
	}
	return users, nil
}

func (ur userRepositoryImple) Update(ctx context.Context, id string, u *user.User) (*user.User, error) {
	documentRef := ur.client.Collection("users").Doc(id)
	_, err := documentRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/user-not-found-on-update",
				Title:      "Error Type Not Found On Update",
				Detail:     "No user found for the provided id when searching to update",
				HTTPStatus: http.StatusNotFound,
				Instance:   "localhost:8080/v1/error/type=instance&identifier=",
				Severity:   problemdetails.Err,

				Err: err,
			})
		}
		return nil, err
	}
	u.Id = id
	_, err = documentRef.Set(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (ur userRepositoryImple) Delete(ctx context.Context, id string) error {
	documentRef := ur.client.Collection("users").Doc(id)
	_, err := documentRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/user-not-found-on-delete",
				Title:      "Error Type Not Found On Delete",
				Detail:     "No user found for the provided id when searching to delete",
				HTTPStatus: http.StatusNotFound,
				Instance:   "localhost:8080/v1/error/type=instance&identifier=",
				Severity:   problemdetails.Err,

				Err: err,
			})
		}
		return err
	}
	_, err = documentRef.Delete(ctx)
	return err
}
