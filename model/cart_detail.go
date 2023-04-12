package model

type CartDetail struct {
	Item     Item `json:"item"`
	Quantity int  `json:"quantity"`
}
