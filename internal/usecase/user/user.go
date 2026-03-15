package user

import "os/user"

type UserUseCase interface {
	Create(*user.User) (*user.User, error)
	Get(string, string) ([]user.User, error)
	Update(string, *user.User) (*user.User, error)
	Delete(string) error
}

type userUseCaseImpl struct {
}

func NewUserUseCase() UserUseCase {
	return userUseCaseImpl{}
}

func (ref userUseCaseImpl) Create(user *user.User) (*user.User, error) {
	return nil, nil
}

func (ref userUseCaseImpl) Get(id string, zipCode string) ([]user.User, error) {
	return nil, nil
}

func (ref userUseCaseImpl) Update(id string, user *user.User) (*user.User, error) {
	return nil, nil
}

func (ref userUseCaseImpl) Delete(id string) error {
	return nil
}
