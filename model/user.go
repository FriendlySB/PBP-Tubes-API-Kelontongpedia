package model

type User struct {
	ID              int         `json:"id"`
	Name            string      `json:"name"`
	Email           string      `json:"email"`
	Password        string      `json:"password,omitempty"`
	Address         string      `json:"address"`
	TelephoneNo     string      `json:"telephone"`
	UserType        int         `json:"usertype,omitempty"`
	BanStatus       bool        `json:"banstatus,omitempty"`
	TransactionList *Transaction `json:"transactionlist,omitempty"`
	Cart            *Cart        `json:"cart,omitempty"`
}
