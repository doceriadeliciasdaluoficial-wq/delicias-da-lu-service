package main

import (
	"context"

	"delicias-da-lu-service.com/mod/internal/controller"
	"delicias-da-lu-service.com/mod/internal/controller/system"
	"delicias-da-lu-service.com/mod/internal/repository/errorFirestore"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func main() {

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "project-4419255d-5de2-41f6-82b")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Firestore client")
	}
	defer client.Close()

	server := controller.NewAPIServer()

	errorRepository := errorFirestore.NewErrorRepository(client)

	testeHandler := system.NewHandler(errorRepository)

	server.AddRoutes(testeHandler)

	server.Start()
}
