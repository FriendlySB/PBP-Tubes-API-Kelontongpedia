package model

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Product `json:"data"`
}
