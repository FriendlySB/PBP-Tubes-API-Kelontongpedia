package model

type Shop struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Reputation    int    `json:"reputation"`
	Category      string `json:"category"`
	Address       string `json:"address"`
	TelephoneNo   string `json:"telephone"`
	Email         string `json:"email"`
	ShopBanStatus bool   `json:"shopbanstatus,omitempty"`
}
