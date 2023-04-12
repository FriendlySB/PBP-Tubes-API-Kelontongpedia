package controller

import (
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
	sqlStatement := `Update users
	SET banstatus = 1
	where ID =?`

	_, errQuery := db.Exec(sqlStatement, userId)

	if errQuery != nil {
		sendErrorResponse(w, "Failed to ban user")
		return
	} else {
		sendSuccessResponse(w, "user banned", nil)
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
	sqlStatement := `Update shop
	SET shopstatus = 1
	where ID =?`

	_, errQuery := db.Exec(sqlStatement, userId)

	if errQuery != nil {
		sendErrorResponse(w, "Failed to ban shop")
		return
	} else {
		sendSuccessResponse(w, "shop banned", nil)
	}
}
