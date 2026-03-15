package main

import (
	"context"

	"delicias-da-lu-service.com/mod/internal/controller"
	"delicias-da-lu-service.com/mod/internal/controller/system"
	userHandlers "delicias-da-lu-service.com/mod/internal/controller/user"
	"delicias-da-lu-service.com/mod/internal/repository/errorFirestore"
	userFirestore "delicias-da-lu-service.com/mod/internal/repository/user"
	"delicias-da-lu-service.com/mod/internal/usecase/errorList"
	userUsecase "delicias-da-lu-service.com/mod/internal/usecase/user"

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
	errorUsecase := errorList.NewErrorListUseCase(errorRepository)
	systemHandler := system.NewHandler(errorUsecase)

	userRepository := userFirestore.NewUserRepository(client)
	userUC := userUsecase.NewUserUseCase(userRepository)
	userHandler := userHandlers.NewUserHandler(userUC)

	server.AddRoutes(systemHandler, userHandler)

	server.Start()
}
