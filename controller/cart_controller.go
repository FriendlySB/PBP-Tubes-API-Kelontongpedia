package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	err := r.ParseForm()
	if err != nil {
		return
	}
	vars := mux.Vars(r)
	UserID := vars["user_id"]

	query := "SELECT a.cartId,b.quantity,c.itemId,c.shopId,c.itemName,c.itemDesc,c.itemCategory,c.itemPrice,c.itemStock FROM cart a INNER JOIN cart_detail b ON a.cartId = b.cartId INNER JOIN item c ON b.itemId = c.itemId WHERE a.userId ='" + UserID + "'"

	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var cart model.Cart
	var cartDetail model.CartDetail
	var cartDetails []model.CartDetail
	var item model.Item
	for rows.Next() {
		if err := rows.Scan(&cart.ID, &cartDetail.Quantity, &item.ID, &item.ShopID, &item.Name, &item.Desc, &item.Category, &item.Price, &item.Stock); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		} else {
			cartDetail.Item = item
			cartDetails = append(cartDetails, cartDetail)
		}
	}
	cart.CartDetail = cartDetails
	var response model.CartResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = cart
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func InsertItemToCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	//Read From Request Body
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed")
		return
	}
	itemId, _ := strconv.Atoi(r.Form.Get("itemId"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	vars := mux.Vars(r)
	UserID := vars["user_id"]

	query := "SELECT cartId FROM cart WHERE userId ='" + UserID + "'"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var cart model.Cart
	for rows.Next() {
		if err := rows.Scan(&cart.ID); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		} else {
			_, errQuery := db.Exec("INSERT INTO cart_detail(cartId,itemId,quantity)values(?,?,?)",
				itemId,
				cart.ID,
				quantity,
			)
			var response model.CartResponse
			if errQuery == nil {
				response.Status = 200
				response.Message = "Insert Item ke Cart Berhasil"
			} else {
				fmt.Println(errQuery)
				response.Status = 400
				response.Message = "Insert Item ke Cart Gagal"
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
		}
	}
}
func UpdateCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
}
func DeleteItemInCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
}
