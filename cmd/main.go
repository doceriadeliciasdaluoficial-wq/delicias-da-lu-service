package main

import (
	"delicias-da-lu-service/internal/controller"
	"delicias-da-lu-service/internal/controller/teste"
)

func main() {
	server := controller.NewAPIServer()

	testeHandler := teste.NewHandler()

	server.AddRoutes(testeHandler)

	server.Start()
}
