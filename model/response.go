package model

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
type CartResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Cart   `json:"cart"`
}
