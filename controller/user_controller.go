package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"log"
	"net/http"
)

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	//Read From Request Body
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed")
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	query := "SELECT userid,Name,Address,UserType FROM USERS WHERE Email ='" + email + "' && Password='" + password + "'"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var user model.User
	login := false
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Address, &user.UserType); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		} else {
			generateToken(w, user.ID, user.Name, user.UserType)
			login = true
		}
	}
	if login {
		sendSuccessResponse(w, "Login Success", user)
	} else {
		sendErrorResponse(w, "Login Failed")
	}
}
func Logout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)
	var user model.User
	sendSuccessResponse(w, "Logout Success", user)
}
