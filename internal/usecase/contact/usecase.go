package contact

import (
	"context"

	"delicias-da-lu-service.com/mod/internal/entity/contact"
	contactRepo "delicias-da-lu-service.com/mod/internal/repository/contact"
)

type ContactUseCase interface {
	Get(ctx context.Context) (*contact.Contact, error)
	Update(ctx context.Context, cnt *contact.Contact) (*contact.Contact, error)
}

type contactUseCaseImpl struct {
	repository contactRepo.ContactRepository
}

func NewContactUseCase(repository contactRepo.ContactRepository) ContactUseCase {
	return contactUseCaseImpl{
		repository: repository,
	}
}

func (c contactUseCaseImpl) Get(ctx context.Context) (*contact.Contact, error) {
	return c.repository.Get(ctx)
}

func (c contactUseCaseImpl) Update(ctx context.Context, cnt *contact.Contact) (*contact.Contact, error) {
	return c.repository.Update(ctx, cnt)
}
