package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
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
		query := "SELECT userid, name, email, address, telpNo FROM users WHERE userid =" + strconv.Itoa(currentID)
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
			if err == nil {
				response.Status = 200
				response.Message = "Success"
				response.Data = users
			} else {
				response.Status = 400
				response.Message = "Error"
			}
			json.NewEncoder(w).Encode(response)
		}
	}
}

func getUserIdFromCookie(r *http.Request) int {
	if cookie, err := r.Cookie(tokenName); err == nil {
		jwtToken := cookie.Value
		accessClaims := &model.Claim{}
		parsedToken, err := jwt.ParseWithClaims(jwtToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return accessClaims.ID
		}
	}
	return -1
}
