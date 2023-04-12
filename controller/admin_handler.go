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
		sendErrorResponse(w, "Failed to update transaction")
		return
	} else {
		sendSuccessResponse(w, "Progress updated", nil)
	}
}
