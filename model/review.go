package model

import "time"

type Review struct {
	ID         int       `json:"id"`
	UserId     int       `json:"userid"`
	ReviewDate time.Time `json:"reviewdate"`
	Rating     int       `json:"rating"`
	Review     string    `json:"review"`
}