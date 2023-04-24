package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func BanUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	vars := mux.Vars(r)
	userId := vars["user_id"]
	sqlStatement := "UPDATE users SET banstatus = 1 WHERE userid =?"

	_, errQuery := db.Exec(sqlStatement, userId)

	if errQuery != nil {
		log.Println(errQuery)
		sendErrorResponse(w, "Failed to ban user")
		return
	} else {
		user := model.User{}
        row := db.QueryRow("SELECT userid, name, email FROM users WHERE userid=?", userId)
        err := row.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
            log.Println(err)
            return
        }
		sendSuccessResponse(w, "User banned", nil)
		sendMailBanUser(user)
	}
}

func BanShop(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	vars := mux.Vars(r)
	userId := vars["shop_id"]
	sqlStatement := "UPDATE shop SET shopstatus = 1 WHERE shopid =?"
	
	_, errQuery := db.Exec(sqlStatement, userId)

	if errQuery != nil {
		log.Println(errQuery)
		sendErrorResponse(w, "Failed to ban shop")
		return
	} else {
		sendSuccessResponse(w, "Shop banned", nil)
	}
}
