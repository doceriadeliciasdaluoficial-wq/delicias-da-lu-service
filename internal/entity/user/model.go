package user

import "delicias-da-lu-service.com/mod/internal/entity/order"

type User struct {
	Id       string `json:"id" firestore:"-"`
	Name     string `json:"name" firestore:"name"`
	Birthday string `json:"birthday" firestore:"birthday"`
	Document string `json:"document" firestore:"document"`
	ZipCode  string `json:"zip_code" firestore:"zip_code" `

	Orders []order.Order `json:"orders"`
}
