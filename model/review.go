package model

import "time"

type Review struct {
	ID         int       `json:"id"`
	ItemID     int       `json:"itemid"`
	ReviewDate time.Time `json:"reviewdate"`
	Rating     int       `json:"rating"`
	Review     string    `json:"review"`
}
