package errorList

import (
	"context"
	"time"

	"delicias-da-lu-service.com/mod/internal/entity/issue"
	"delicias-da-lu-service.com/mod/internal/repository/errorFirestore"
)

type CreateErrorTypeRequest struct {
	Identifier string `json:"identifier"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	Status     int    `json:"status"`
}

type ErrorListUseCase interface {
	GetTypeOfErrorByIdentifier(context.Context, string) (string, error)
	GetInstanceOfErrorByIdentifier(context.Context, string) (issue.ErrorInstance, error)
	GetInstanceHTMLByIdentifier(context.Context, string) (string, error)
	SeedErrorTypes(context.Context) error
	CreateErrorTypesFromList(context.Context, []CreateErrorTypeRequest) error
	DeleteErrorType(context.Context, string) error
}

type errorListUseCaseImpl struct {
	errorRepository errorFirestore.ErrorRepository
}

func NewErrorListUseCase(repository errorFirestore.ErrorRepository) ErrorListUseCase {
	return errorListUseCaseImpl{
		errorRepository: repository,
	}
}

func (ref errorListUseCaseImpl) GetTypeOfErrorByIdentifier(ctx context.Context, identifier string) (string, error) {
	content, err := ref.errorRepository.GetTypeOfErrorByIdentifier(ctx, identifier)
	if err != nil {
		return "", err
	}
	return content.Html, err
}

func (ref errorListUseCaseImpl) GetInstanceOfErrorByIdentifier(ctx context.Context, identifier string) (issue.ErrorInstance, error) {
	return ref.errorRepository.GetInstanceOfErrorByIdentifier(ctx, identifier)
}

func (ref errorListUseCaseImpl) GetInstanceHTMLByIdentifier(ctx context.Context, identifier string) (string, error) {
	instance, err := ref.errorRepository.GetInstanceOfErrorByIdentifier(ctx, identifier)
	if err != nil {
		return "", err
	}
	return instance.Html, nil
}

func (ref errorListUseCaseImpl) CreateErrorTypesFromList(ctx context.Context, requests []CreateErrorTypeRequest) error {
	for _, req := range requests {
		html := buildTypeHTML(typeHTMLData{
			Identifier: req.Identifier,
			Title:      req.Title,
			Detail:     req.Detail,
			Status:     req.Status,
			UpdatedAt:  time.Now().UTC(),
		})

		err := ref.errorRepository.UpsertErrorType(ctx, req.Identifier, issue.ErrorType{
			Html:      html,
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (ref errorListUseCaseImpl) DeleteErrorType(ctx context.Context, identifier string) error {
	return ref.errorRepository.DeleteErrorType(ctx, identifier)
}
