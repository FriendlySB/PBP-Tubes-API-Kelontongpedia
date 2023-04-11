package model

type Cart struct {
	ID         int        `json:"id"`
	CartDetail CartDetail `json:"cartdetail"`
}
