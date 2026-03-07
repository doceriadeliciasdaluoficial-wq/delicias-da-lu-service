package main

import (
	"delicias-da-lu-service/internal/controller"
	"delicias-da-lu-service/internal/controller/system"
)

func main() {
	server := controller.NewAPIServer()

	testeHandler := system.NewHandler()

	server.AddRoutes(testeHandler)

	server.Start()
}
