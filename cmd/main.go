package main

import (
	"context"
	"os"

	"delicias-da-lu-service.com/mod/internal/controller"
	authController "delicias-da-lu-service.com/mod/internal/controller/auth"
	cakeBuilderController "delicias-da-lu-service.com/mod/internal/controller/cakebuilder"
	configController "delicias-da-lu-service.com/mod/internal/controller/config"
	contactController "delicias-da-lu-service.com/mod/internal/controller/contact"
	healthController "delicias-da-lu-service.com/mod/internal/controller/health"
	homeController "delicias-da-lu-service.com/mod/internal/controller/home"
	menuController "delicias-da-lu-service.com/mod/internal/controller/menu"
	orderController "delicias-da-lu-service.com/mod/internal/controller/order"
	systemController "delicias-da-lu-service.com/mod/internal/controller/system"
	uploadController "delicias-da-lu-service.com/mod/internal/controller/upload"
	"delicias-da-lu-service.com/mod/internal/platform/problemdetails"
	repoCakeBuilder "delicias-da-lu-service.com/mod/internal/repository/cakebuilder"
	repoConfig "delicias-da-lu-service.com/mod/internal/repository/config"
	repoContact "delicias-da-lu-service.com/mod/internal/repository/contact"
	repoError "delicias-da-lu-service.com/mod/internal/repository/errorFirestore"
	repoHome "delicias-da-lu-service.com/mod/internal/repository/home"
	repoMenu "delicias-da-lu-service.com/mod/internal/repository/menu"
	repoOrder "delicias-da-lu-service.com/mod/internal/repository/order"
	repoStorage "delicias-da-lu-service.com/mod/internal/repository/storage"
	repoUser "delicias-da-lu-service.com/mod/internal/repository/user"
	"delicias-da-lu-service.com/mod/internal/usecase/auth"
	"delicias-da-lu-service.com/mod/internal/usecase/cakebuilder"
	"delicias-da-lu-service.com/mod/internal/usecase/config"
	"delicias-da-lu-service.com/mod/internal/usecase/contact"
	"delicias-da-lu-service.com/mod/internal/usecase/errorList"
	homeUsecase "delicias-da-lu-service.com/mod/internal/usecase/home"
	"delicias-da-lu-service.com/mod/internal/usecase/menu"
	"delicias-da-lu-service.com/mod/internal/usecase/order"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"
)

func main() {
	logger := log.With().Str("component", "bootstrap").Logger()
	ctx := context.Background()

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = "project-4419255d-5de2-41f6-82b"
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		logger.Fatal().Err(err).Str("project_id", projectID).Msg("failed to create Firestore client")
	}
	defer client.Close()

	logger.Info().Str("project_id", projectID).Msg("firestore client initialized")

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create Cloud Storage client")
	}
	defer storageClient.Close()

	logger.Info().Msg("cloud storage client initialized")

	userRepository := repoUser.NewUserRepository(client)
	menuRepository := repoMenu.NewMenuRepository(client)
	cakeBuilderRepository := repoCakeBuilder.NewCakeBuilderRepository(client)
	contactRepository := repoContact.NewContactRepository(client)
	orderRepository := repoOrder.NewOrderRepository(client)
	configRepository := repoConfig.NewConfigRepository(client)
	errorRepository := repoError.NewErrorRepository(client)
	homeRepository := repoHome.NewHomeRepository(client)
	storageRepository := repoStorage.NewStorageRepository(storageClient)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	authUseCase := auth.NewAuthUseCase(userRepository, jwtSecret)
	menuUseCase := menu.NewMenuUseCase(menuRepository)
	cakeBuilderUseCase := cakebuilder.NewCakeBuilderUseCase(cakeBuilderRepository)
	contactUseCase := contact.NewContactUseCase(contactRepository)
	orderUseCase := order.NewOrderUseCase(orderRepository)
	configUseCase := config.NewConfigUseCase(configRepository)
	errorListUseCase := errorList.NewErrorListUseCase(errorRepository)
	homeUseCase := homeUsecase.NewHomeUseCase(homeRepository)
	errorRecorder := errorList.NewErrorRecorder(errorRepository)

	problemdetails.SetErrorRecorder(errorRecorder)
	if err := errorListUseCase.SeedErrorTypes(ctx); err != nil {
		logger.Warn().Err(err).Msg("failed to seed error types")
	}

	healthHandler := healthController.NewHealthHandler()
	authHandler := authController.NewAuthHandler(authUseCase)
	menuHandler := menuController.NewMenuHandler(menuUseCase)
	cakeBuilderHandler := cakeBuilderController.NewCakeBuilderHandler(cakeBuilderUseCase)
	contactHandler := contactController.NewContactHandler(contactUseCase)
	orderHandler := orderController.NewOrderHandler(orderUseCase)
	configHandler := configController.NewConfigHandler(configUseCase)
	systemHandler := systemController.NewHandler(errorListUseCase)
	homeHandler := homeController.NewHomeHandler(homeUseCase)
	uploadHandler := uploadController.NewUploadHandler(storageRepository)

	// Initialize API server
	server := controller.NewAPIServer()

	// Register routes
	if err := server.RegisterRoutes(
		healthHandler,
		authHandler,
		authUseCase,
		menuHandler,
		cakeBuilderHandler,
		contactHandler,
		configHandler,
		orderHandler,
		systemHandler,
		homeHandler,
		uploadHandler,
	); err != nil {
		logger.Fatal().Err(err).Msg("failed to register routes")
	}

	logger.Info().Msg("routes registered")

	// Start server
	if err := server.Start(); err != nil {
		logger.Fatal().Err(err).Msg("server failed to start")
	}
}
