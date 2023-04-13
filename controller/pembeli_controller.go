package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"encoding/json"
	"log"
	"net/http"
)

// "encoding/json"
// "fmt"

// "github.com/gorilla/mux"

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	currentID := getUserIdFromCookie(r)

	if currentID == -1 {
		sendUnauthorizedResponse(w)
	} else {
		query := "SELECT userid, name, email, address, telpNo FROM users"
		name := r.URL.Query()["name"]
		userid := r.URL.Query()["userid"]
		if name != nil {
			query += " WHERE name='" + name[0] + "'"
		}
		if userid != nil {
			query += " WHERE userid=" + userid[0]
		}

		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		} else {
			var user model.User
			var users []model.User
			for rows.Next() {
				if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.TelephoneNo); err != nil {
					sendErrorResponse(w, "Error while scanning rows")
					return
				} else {
					users = append(users, user)
				}
			}
			var response model.GenericResponse
			w.Header().Set("Content=Type", "application/json")
			if err == nil {
				response.Status = 200
				response.Message = "Success"
				response.Data = users
				json.NewEncoder(w).Encode(response)
			} else {
				response.Status = 400
				response.Message = "Error"
				json.NewEncoder(w).Encode(response)
			}
		}
	}
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	currentID := getUserIdFromCookie(r)

	if currentID == -1 {
		sendUnauthorizedResponse(w)
	} else {
		err := r.ParseForm()
		if err != nil {
			sendErrorResponse(w, "Error while parsing form")
			return
		}
		name := r.Form.Get("name")
		email := r.Form.Get("email")
		address := r.Form.Get("address")
		telpNo := r.Form.Get("telpNo")
		
		_, errQuery := db.Exec("UPDATE users SET name=?, email=?, address=?, telpNo=? WHERE userid=?", name, email, address, telpNo, currentID)

		var response model.GenericResponse
		if errQuery == nil {
			response.Status = 200
			response.Message = "Success"
			json.NewEncoder(w).Encode(response)
		} else {
			response.Status = 400
			response.Message = "Error"
			json.NewEncoder(w).Encode(response)
		}
	}
}

func ReviewItem(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	currentID := getUserIdFromCookie(r)

	if currentID == -1 {
		sendUnauthorizedResponse(w)
	}
}