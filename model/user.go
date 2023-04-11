package model

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	TelephoneNo string `json:"telephone"`
	UserType    int    `json:"usertype"`
	BanStatus   bool   `json:"banstatus"`
}
