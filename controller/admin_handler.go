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
	shopId := vars["shop_id"]
	sqlStatement := "UPDATE shop SET shopstatus = 1 WHERE shopid =?"

	_, errQuery := db.Exec(sqlStatement, shopId)

	sqlStatement2 := "select a.UserId, a.Name, a.email FROM users a INNER JOIN shop_admin b on a.userid = b.userId where b.shopId = ?"

	rows, err := db.Query(sqlStatement2, shopId)
	if err != nil {
		log.Println(err)
		return
	}

	var user model.User
	var users []model.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			log.Println(err)
			return
		} else {
			users = append(users, user)
		}
	}
	for i := 0; i < len(users); i++ {
		sqlStatement3 := "UPDATE users SET banstatus = 1 WHERE userid =?"
		_, errQuery := db.Exec(sqlStatement3, users[i].ID)
		if errQuery != nil {
			log.Println(errQuery)
			sendErrorResponse(w, "Failed to ban user")
			return
		} else {
			sendSuccessResponse(w, "User banned", nil)
			sendMailBanUser(users[i])
		}
	}

	if errQuery != nil {
		log.Println(errQuery)
		sendErrorResponse(w, "Failed to ban shop")
		return
	} else {
		sendSuccessResponse(w, "Shop banned", nil)
	}
}
