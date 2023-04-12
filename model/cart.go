package model

type Cart struct {
	ID         int          `json:"id"`
	CartDetail []CartDetail `json:"cartdetail"`
}
type UpdateCart struct {
	CartID   int `json:"cartid"`
	ItemID   int `json:"itemid"`
	Quantity int `json:"quantity"`
}
type DeleteCart struct {
	CartID int `json:"cartid"`
	ItemID int `json:"itemid"`
}
