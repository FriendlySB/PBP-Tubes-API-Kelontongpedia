package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("joensuu")
var tokenName = "token"

func generateToken(w http.ResponseWriter, id int, name string, usertype int) {
	tokenExpiryTime := time.Now().Add(5 * time.Minute)

	claims := &model.Claim{
		ID:       id,
		Name:     name,
		UserType: usertype,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err := token.SignedString(jwtKey)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    jwtToken,
		Expires:  tokenExpiryTime,
		Secure:   false,
		HttpOnly: true,
	})
}

func resetUserToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
	})
}

func Authenticate(next http.HandlerFunc, accesstype int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValidToken := validateUserToken(r, accesstype)
		if !isValidToken {
			//sendUnauthorizedResponse(w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func validateUserToken(r *http.Request, accessType int) bool {
	isAccessTokenValid, id, name, usertype := validateTokenFromCookies(r)
	fmt.Println(id, name, usertype, accessType, isAccessTokenValid)

	if isAccessTokenValid {
		isUserValid := usertype == accessType
		if isUserValid {
			return true
		}
	}
	return false
}
func validateTokenFromCookies(r *http.Request) (bool, int, string, int) {
	if cookie, err := r.Cookie(tokenName); err == nil {
		jwtToken := cookie.Value
		accessClaims := &model.Claim{}
		parsedToken, err := jwt.ParseWithClaims(jwtToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.ID, accessClaims.Name, accessClaims.UserType
		}
	}
	return false, -1, "", -1
}

func sendUnauthorizedResponse(w http.ResponseWriter) {
	var response model.ErrorResponse
	response.Status = 401
	response.Message = "Unauthorized Access"
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
