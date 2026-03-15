package user

import (
	"context"

	"delicias-da-lu-service.com/mod/internal/entity/user"
	"delicias-da-lu-service.com/mod/internal/platform/validator"
	userRepository "delicias-da-lu-service.com/mod/internal/repository/user"
)

type UserUseCase interface {
	Create(context.Context, *user.User) (*user.User, error)
	Get(context.Context, string, string) ([]user.User, error)
	Update(context.Context, string, *user.User) (*user.User, error)
	Delete(context.Context, string) error
}

type userUseCaseImpl struct {
	userValidator  validator.Validator[user.User]
	userRepository userRepository.UserRepository
}

func NewUserUseCase(userRepository userRepository.UserRepository) UserUseCase {
	return userUseCaseImpl{
		userValidator:  NewUserValidator(),
		userRepository: userRepository,
	}
}

func (ref userUseCaseImpl) Create(ctx context.Context, user *user.User) (*user.User, error) {
	if err := ref.userValidator.Validate(*user); err != nil {
		return nil, err
	}
	return ref.userRepository.Create(context.Background(), user)
}

func (ref userUseCaseImpl) Get(ctx context.Context, field string, value string) ([]user.User, error) {
	if field == "id" {
		field = "DocumentId"
	}
	return ref.userRepository.Get(ctx, field, value)
}

func (ref userUseCaseImpl) Update(ctx context.Context, id string, user *user.User) (*user.User, error) {
	if err := ref.userValidator.Validate(*user); err != nil {
		return nil, err
	}
	return ref.userRepository.Update(ctx, id, user)
}

func (ref userUseCaseImpl) Delete(ctx context.Context, id string) error {
	return ref.userRepository.Delete(ctx, id)
}
