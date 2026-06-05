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
	menuController "delicias-da-lu-service.com/mod/internal/controller/menu"
	orderController "delicias-da-lu-service.com/mod/internal/controller/order"
	repoCakeBuilder "delicias-da-lu-service.com/mod/internal/repository/cakebuilder"
	repoConfig "delicias-da-lu-service.com/mod/internal/repository/config"
	repoContact "delicias-da-lu-service.com/mod/internal/repository/contact"
	repoMenu "delicias-da-lu-service.com/mod/internal/repository/menu"
	repoOrder "delicias-da-lu-service.com/mod/internal/repository/order"
	repoUser "delicias-da-lu-service.com/mod/internal/repository/user"
	"delicias-da-lu-service.com/mod/internal/usecase/auth"
	"delicias-da-lu-service.com/mod/internal/usecase/cakebuilder"
	"delicias-da-lu-service.com/mod/internal/usecase/config"
	"delicias-da-lu-service.com/mod/internal/usecase/contact"
	"delicias-da-lu-service.com/mod/internal/usecase/menu"
	"delicias-da-lu-service.com/mod/internal/usecase/order"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func main() {

	ctx := context.Background()

	// Initialize Firestore client
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = "project-4419255d-5de2-41f6-82b"
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Firestore client")
	}
	defer client.Close()

	// Initialize repositories
	userRepository := repoUser.NewUserRepository(client)
	menuRepository := repoMenu.NewMenuRepository(client)
	cakeBuilderRepository := repoCakeBuilder.NewCakeBuilderRepository(client)
	contactRepository := repoContact.NewContactRepository(client)
	orderRepository := repoOrder.NewOrderRepository(client)
	configRepository := repoConfig.NewConfigRepository(client)

	// Initialize use cases
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production" // Change in production!
	}

	authUseCase := auth.NewAuthUseCase(userRepository, jwtSecret)
	menuUseCase := menu.NewMenuUseCase(menuRepository)
	cakeBuilderUseCase := cakebuilder.NewCakeBuilderUseCase(cakeBuilderRepository)
	contactUseCase := contact.NewContactUseCase(contactRepository)
	orderUseCase := order.NewOrderUseCase(orderRepository)
	configUseCase := config.NewConfigUseCase(configRepository)

	// Initialize handlers
	healthHandler := healthController.NewHealthHandler()
	authHandler := authController.NewAuthHandler(authUseCase)
	menuHandler := menuController.NewMenuHandler(menuUseCase)
	cakeBuilderHandler := cakeBuilderController.NewCakeBuilderHandler(cakeBuilderUseCase)
	contactHandler := contactController.NewContactHandler(contactUseCase)
	orderHandler := orderController.NewOrderHandler(orderUseCase)
	configHandler := configController.NewConfigHandler(configUseCase)

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
	); err != nil {
		log.Fatal().Err(err).Msg("Failed to register routes")
	}

	// Start server
	if err := server.Start(); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
