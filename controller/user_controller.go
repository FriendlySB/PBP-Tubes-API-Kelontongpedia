package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"encoding/json"
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
		sendSuccessResponse(w, "Login Success", nil)
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
	sendSuccessResponse(w, "Logout Success", nil)
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
			sendErrorResponse(w, "Error")
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

		if errQuery == nil {
			rows, err := db.Query("SELECT userid, name, email, address, telpNo FROM users WHERE userid = ?", currentID)
			if err != nil {
				sendErrorResponse(w, "Error while fetching updated data")
				return
			}
			defer rows.Close()

			var user model.User
			for rows.Next() {
				if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.TelephoneNo); err != nil {
					sendErrorResponse(w, "Error while scanning rows")
					return
				}
			}
			response := model.GenericResponse{Status: 200, Message: "Success", Data: user}
			json.NewEncoder(w).Encode(response)
		} else {
			response := model.GenericResponse{Status: 400, Message: "Error"}
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
	} else {
		err := r.ParseForm()
		if err != nil {
			sendErrorResponse(w, "Error while parsing form")
			return
		}
		itemId := r.Form.Get("itemId")
		rating := r.Form.Get("rating")
		review := r.Form.Get("review")

		_, errQuery := db.Exec("INSERT INTO review(itemID, userID, rating, review) VALUES(?,?,?,?)", itemId, currentID, rating, review)
		
		if errQuery == nil {
			rows, err := db.Query("SELECT reviewId, userId, review_date, rating, review FROM review WHERE reviewId = LAST_INSERT_ID()")
			if err != nil {
				sendErrorResponse(w, "Error while fetching updated data")
				return
			}
			defer rows.Close()

			var review model.Review
			for rows.Next() {
				if err := rows.Scan(&review.ID, &review.UserId, &review.ReviewDate, &review.Rating, &review.Review); err != nil {
					sendErrorResponse(w, "Error while scanning rows")
					return
				}
			}
			response := model.GenericResponse{Status: 200, Message: "Success", Data: review}
			json.NewEncoder(w).Encode(response)
		} else {
			response := model.GenericResponse{Status: 400, Message: "Error"}
			json.NewEncoder(w).Encode(response)
		}
	}
}