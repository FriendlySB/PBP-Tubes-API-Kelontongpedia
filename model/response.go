package model

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
type GenericResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
type CartResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Cart   `json:"cart"`
}
type ReviewResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    []Review `json:"review"`
}
