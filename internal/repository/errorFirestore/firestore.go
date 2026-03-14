package errorFirestore

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/firestore"
)

type ErrorRepository interface {
	GetTypeOfErrorByIdentifier(ctx context.Context, identifier string) (string, error)
	GetInstanceOfErrorByIdentifier(ctx context.Context, identifier string) (string, error)
}

type errorRepositoryImple struct {
	client *firestore.Client
}

func NewErrorRepository(client *firestore.Client) ErrorRepository {
	return errorRepositoryImple{
		client: client,
	}
}

func (ref errorRepositoryImple) GetTypeOfErrorByIdentifier(ctx context.Context, identifier string) (string, error) {
	doc, err := ref.client.Collection("type").Doc(identifier).Get(ctx)
	if err != nil {
		return "", err
	}

	type errorTypeDAO struct {
		html string `firestore:"html"`
	}

	var errorType errorTypeDAO
	if err := doc.DataTo(&errorType); err != nil {
		return "", err
	}
	return errorType.html, nil
}
func (ref errorRepositoryImple) GetInstanceOfErrorByIdentifier(ctx context.Context, identifier string) (string, error) {
	doc, err := ref.client.Collection("instance").Doc(identifier).Get(ctx)
	if err != nil {
		return "", err
	}

	type errorInstanceDAO struct {
		RequestBody string `firestore:"request_body"`
		RequestDate string `firestore:"request_date"`
		Status      int    `firestore:"status"`
		Title       string `firestore:"title"`
		TraceID     string `firestore:"trace_id"`
		Type        string `firestore:"type"`
		UserAgent   string `firestore:"user_agent"`
	}

	var errorInstance errorInstanceDAO
	if err := doc.DataTo(&errorInstance); err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(errorInstance)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
