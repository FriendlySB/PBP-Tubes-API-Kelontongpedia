package model

type TransactionDetail struct {
	IdTransaction int  `json:"idtransaction"`
	Item          Item `json:"item"`
	Quantity      int  `json:"quantity"`
}
