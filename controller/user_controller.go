package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	var user model.User
	_, UserID, name, _ := validateTokenFromCookies(r)
	user.ID = UserID
	user.Name = name
	resetUserToken(w)
	sendSuccessResponse(w, "Logout Success", user)
}

// fungsi untuk register user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	//Read From Request Body
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Failed")
		return
	}
	name := r.Form.Get("name")
	address := r.Form.Get("address")
	res, errQuery := db.Exec("INSERT INTO users(name,address)values(?,?,?)",
		name,
		address,
	)
	id, _ := res.LastInsertId()
	var user model.User
	user.ID = int(id)
	user.Name = name
	user.Address = address
	if errQuery == nil {
		sendSuccessResponse(w, "Register Berhasil", user)
	} else {
		fmt.Println(errQuery)
		sendErrorResponse(w, "Register Gagal")
	}
}
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	//Read From Request Body
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed")
		return
	}
	//Ambil password lama dan baru dari form
	oldpassword := r.Form.Get("old_password")
	newpassword := r.Form.Get("new_password")
	//Password lama user dari database untuk dibandingkan
	var password string

	//User id ambil pakai cookie
	userid := 3

	query := "SELECT password FROM users WHERE userid = " + strconv.Itoa(userid)
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	for rows.Next() {
		if err := rows.Scan(&password); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		}
	}
	if password == oldpassword {
		query = "UPDATE users SET password = '" + newpassword + "' WHERE userid = " + strconv.Itoa(userid)
		_, errQuery := db.Exec(query)

		if errQuery != nil {
			sendErrorResponse(w, "Failed to change password!")
		} else {
			sendSuccessResponse(w, "Password successfully changed!", nil)
		}
	} else {
		sendErrorResponse(w, "Password does not match!")
	}
}
