package model

type Transaction struct {
	ID        int `json:"id"`
	UserID    int `json:"userID"`
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}

type DetailedTransaction struct {
	ID       int     `json:"id"`
	User     User    `json:"user"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type TransactionResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []Transaction `json:"data"`
}

type DetailedTransactionResponse struct {
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Data    []DetailedTransaction `json:"data"`
}
