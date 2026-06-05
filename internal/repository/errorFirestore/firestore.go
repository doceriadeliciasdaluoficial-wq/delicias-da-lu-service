package errorFirestore

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"delicias-da-lu-service.com/mod/internal/entity/issue"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorRepository interface {
	GetTypeOfErrorByIdentifier(ctx context.Context, identifier string) (issue.ErrorType, error)
	GetInstanceOfErrorByIdentifier(ctx context.Context, identifier string) (issue.ErrorInstance, error)
	UpsertErrorType(ctx context.Context, identifier string, errorType issue.ErrorType) error
	CreateErrorInstance(ctx context.Context, identifier string, errorInstance issue.ErrorInstance) error
}

type errorRepositoryImple struct {
	client *firestore.Client
}

func NewErrorRepository(client *firestore.Client) ErrorRepository {
	return errorRepositoryImple{
		client: client,
	}
}

func (ref errorRepositoryImple) GetTypeOfErrorByIdentifier(ctx context.Context, identifier string) (issue.ErrorType, error) {
	doc, err := ref.client.Collection("types").Doc(identifier).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return issue.ErrorType{}, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/type-not-found",
				Title:      "Error Type Not Found",
				Detail:     "No error type found for the provided identifier",
				HTTPStatus: http.StatusNotFound,
				Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/error/type=instance&identifier=",
				Severity:   problemdetails.Err,
			})
		}
		return issue.ErrorType{}, err
	}

	var errorType issue.ErrorType
	if err := doc.DataTo(&errorType); err != nil {
		return issue.ErrorType{}, err
	}

	if errorType.Html == "" {
		return issue.ErrorType{}, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Title:      "Error Type Not Found",
			Detail:     "No error type found for the provided identifier",
			HTTPStatus: http.StatusNotFound,
			Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/error/type/",
			Severity:   problemdetails.Err,
		})
	}
	return errorType, nil
}
func (ref errorRepositoryImple) GetInstanceOfErrorByIdentifier(ctx context.Context, identifier string) (issue.ErrorInstance, error) {
	doc, err := ref.client.Collection("instances").Doc(identifier).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return issue.ErrorInstance{}, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/instance-not-found",
				Title:      "Error Instance Not Found",
				Detail:     "No error instance found for the provided identifier",
				HTTPStatus: http.StatusNotFound,
				Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/error/type=instance&identifier=",
				Severity:   problemdetails.Err,
			})
		}
		return issue.ErrorInstance{}, err
	}

	if status.Code(err) == codes.NotFound {
		return issue.ErrorInstance{}, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/instance-not-found",
			Title:      "Error Instance Not Found",
			Detail:     "No error instance found for the provided identifier",
			HTTPStatus: http.StatusNotFound,
			Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/error/type=instance&identifier=",
			Severity:   problemdetails.Err,
		})
	}

	var errorInstance issue.ErrorInstance
	if err := doc.DataTo(&errorInstance); err != nil {
		return issue.ErrorInstance{}, err
	}

	return errorInstance, nil
}

func (ref errorRepositoryImple) UpsertErrorType(ctx context.Context, identifier string, errorType issue.ErrorType) error {
	docRef := ref.client.Collection("types").Doc(identifier)

	snapshot, err := docRef.Get(ctx)
	if err == nil {
		var existing issue.ErrorType
		if err := snapshot.DataTo(&existing); err == nil && existing.Html != "" {
			return nil
		}
	} else if status.Code(err) != codes.NotFound {
		return err
	}

	_, err = docRef.Set(ctx, errorType, firestore.MergeAll)
	return err
}

func (ref errorRepositoryImple) CreateErrorInstance(ctx context.Context, identifier string, errorInstance issue.ErrorInstance) error {
	_, err := ref.client.Collection("instances").Doc(identifier).Set(ctx, errorInstance)
	return err
}
