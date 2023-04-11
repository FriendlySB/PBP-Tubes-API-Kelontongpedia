package model

type TransactionDetail struct {
	Item     Item `json:"item"`
	Quantity int  `json:"quantity"`
}
