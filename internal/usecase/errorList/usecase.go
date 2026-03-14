package errorList

import (
	"context"

	"delicias-da-lu-service.com/mod/internal/entity/issue"
	"delicias-da-lu-service.com/mod/internal/repository/errorFirestore"
)

type ErrorListUseCase interface {
	GetTypeOfErrorByIdentifier(context.Context, string) (issue.ErrorType, error)
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

func (ref errorListUseCaseImpl) GetTypeOfErrorByIdentifier(ctx context.Context, identifier string) (issue.ErrorType, error) {
	return ref.errorRepository.GetTypeOfErrorByIdentifier(ctx, identifier)
}

func (ref errorListUseCaseImpl) GetInstanceOfErrorByIdentifier(ctx context.Context, identifier string) (issue.ErrorInstance, error) {
	return ref.errorRepository.GetInstanceOfErrorByIdentifier(ctx, identifier)
}
