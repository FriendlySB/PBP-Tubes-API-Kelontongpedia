package model

type Item struct {
	ID       int    `json:"id"`
	ShopID string `json:"shopid"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Category string `json:"category"`
	Price    int    `json:"price"`
	Stock    int    `json:"stock"`
}
