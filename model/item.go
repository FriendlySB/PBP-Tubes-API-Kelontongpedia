package model

type Item struct {
	ID       int    `json:"id"`
	ShopID   string `json:"shopid,omitempty"`
	Name     string `json:"name,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Category string `json:"category,omitempty"`
	Price    int    `json:"price,omitempty"`
	Stock    int    `json:"stock,omitempty"`
}
