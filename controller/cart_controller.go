package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"encoding/json"
	"log"
	"net/http"

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

	query := "SELECT a.cartId, b.cartId,b.itemId,b.quantity,c.itemId,c.shopId,c.itemName,c.itemDesc,c.itemCategory,c.itemPrice,c.itemStock FROM cart a INNER JOIN cart_detail b ON a.cartId = b.cartId INNER JOIN item c ON b.itemId = c.itemId WHERE a.userId ='" + UserID + "'"

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
	var items []model.Item
	for rows.Next() {
		if err := rows.Scan(&cart.ID, &cartDetail.Quantity, &item.ID, &item.ShopID, &item.Name, &item.Desc, &item.Category, &item.Price, &item.Stock); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		} else {
			items = append(items, item)
			cartDetail.Item = items
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
}
func UpdateCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
}
func DeleteItemInCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
}
