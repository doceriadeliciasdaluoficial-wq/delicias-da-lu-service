package contact

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"delicias-da-lu-service.com/mod/internal/entity/contact"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const contactDocID = "contact"

type ContactRepository interface {
	Get(ctx context.Context) (*contact.Contact, error)
	Update(ctx context.Context, cnt *contact.Contact) (*contact.Contact, error)
}

type contactRepositoryImpl struct {
	client *firestore.Client
}

func NewContactRepository(client *firestore.Client) ContactRepository {
	return contactRepositoryImpl{
		client: client,
	}
}

func (r contactRepositoryImpl) Get(ctx context.Context) (*contact.Contact, error) {
	doc, err := r.client.Collection("contacts").Doc(contactDocID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
				Type:       "https://delicias-da-lu-service.com/docs/errors/not-found",
				Title:      "Contacts Not Found",
				Detail:     "No contact information available. Please configure contacts in the admin panel.",
				HTTPStatus: http.StatusNotFound,
				Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/contacts",
				Severity:   problemdetails.Err,
			})
		}
		return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/database-error",
			Title:      "Database Error",
			Detail:     "Failed to retrieve contacts: " + err.Error(),
			HTTPStatus: http.StatusInternalServerError,
			Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/contacts",
			Severity:   problemdetails.Err,
		})
	}

	var cnt contact.Contact
	if err := doc.DataTo(&cnt); err != nil {
		return nil, problemdetails.NewErrorWithStackTrace(problemdetails.Error{
			Type:       "https://delicias-da-lu-service.com/docs/errors/parse-error",
			Title:      "Contact Parse Error",
			Detail:     "Failed to parse contact data: " + err.Error(),
			HTTPStatus: http.StatusInternalServerError,
			Instance:   "https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/contacts",
			Severity:   problemdetails.Err,
		})
	}

	return &cnt, nil
}

func (r contactRepositoryImpl) Update(ctx context.Context, cnt *contact.Contact) (*contact.Contact, error) {
	_, err := r.client.Collection("contacts").Doc(contactDocID).Set(ctx, cnt)
	if err != nil {
		return nil, err
	}

	return cnt, nil
}
