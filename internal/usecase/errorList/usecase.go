package errorList

import (
	"context"

	"delicias-da-lu-service.com/mod/internal/entity/issue"
	"delicias-da-lu-service.com/mod/internal/repository/errorFirestore"
)

type ErrorListUseCase interface {
	GetTypeOfErrorByIdentifier(context.Context, string) (string, error)
	GetInstanceOfErrorByIdentifier(context.Context, string) (issue.ErrorInstance, error)
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
