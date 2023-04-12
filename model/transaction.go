package model

import "time"

type Transaction struct {
	ID                int                 `json:"id"`
	Address           string              `json:"address"`
	Date              time.Time           `json:"date"`
	Delivery          string              `json:"delivery"`
	Progress          string              `json:"progress"`
	PaymentType       string              `json:"paymenttype"`
	TransactionDetail []TransactionDetail `json:"transactiondetail"`
}
type UpdateTransaction struct {
	TransactionID int `json:"transactionid"`
	ItemID        int `json:"itemid"`
	Quantity      int `json:"quantity"`
}
