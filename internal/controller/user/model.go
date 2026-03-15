package user

import "delicias-da-lu-service.com/mod/internal/controller/order"

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Document string `json:"document"`
	ZipCode  string `json:"zip_code" `

	Orders []order.Order `json:"orders"`
}
