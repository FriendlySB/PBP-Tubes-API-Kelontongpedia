package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// get cart dari user yang sedang berjalan
func GetCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	err := r.ParseForm()
	if err != nil {
		return
	}
	// ambil iduser dari cookie
	_, UserID, _, _ := validateTokenFromCookies(r)
	//jika user belum login, maka akan direturn unauthorized response
	if UserID == -1 {
		sendUnauthorizedResponse(w)
		return
	}
	//query ke database
	query := "SELECT a.cartId,b.quantity,c.itemId,c.shopId,c.itemName,c.itemDesc,c.itemCategory,c.itemPrice,c.itemStock FROM cart a INNER JOIN cart_detail b ON a.cartId = b.cartId INNER JOIN item c ON b.itemId = c.itemId WHERE a.userId ='" + strconv.Itoa(UserID) + "'"
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
	sendSuccessResponse(w, "Get Item di Cart Berhasil", cart)
}

// insert item ke cart... asumsi insert itemnya itu item yang tidak ada di cart, kalau itemnya ada, berarti pakai update
func InsertItemToCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed")
		return
	}
	//baca dari request body
	itemId, _ := strconv.Atoi(r.Form.Get("itemId"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))
	//ambil iduser dari cookie
	_, UserID, _, _ := validateTokenFromCookies(r)
	//jika user belum login, maka akan direturn unauthorized response
	if UserID == -1 {
		sendUnauthorizedResponse(w)
		return
	}
	query := "SELECT cartId FROM cart WHERE userId ='" + strconv.Itoa(UserID) + "'"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var cart model.Cart
	var cartDetail model.CartDetail
	for rows.Next() {
		if err := rows.Scan(&cart.ID); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		} else {
			_, errQuery := db.Exec("INSERT INTO cart_detail(cartId,itemId,quantity)values(?,?,?)",
				cart.ID,
				itemId,
				quantity,
			)
			cartDetail.Item.ID = itemId
			cartDetail.Quantity = quantity
			cart.CartDetail = append(cart.CartDetail, cartDetail)
			if errQuery == nil {
				sendSuccessResponse(w, "Insert Item ke Cart Berhasil", cart)
			} else {
				log.Println(errQuery)
				sendErrorResponse(w, "Insert Item ke cart gagal")
			}
		}
	}
}

// update quantity cart,  asumsi item nya selalu sudah ada di cart
func UpdateCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	//Read From Request Body
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Failed")
		return
	}
	itemId, _ := strconv.Atoi(r.Form.Get("itemId"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	_, UserID, _, _ := validateTokenFromCookies(r)
	//jika user belum login, maka akan direturn unauthorized response
	if UserID == -1 {
		sendUnauthorizedResponse(w)
		return
	}
	query := "SELECT cartId FROM cart WHERE userId ='" + strconv.Itoa(UserID) + "'"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var cart model.Cart
	var cartDetail model.CartDetail
	for rows.Next() {
		if err := rows.Scan(&cart.ID); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		} else {
			//belum cek item nya ada atau gk di cart
			_, errQuery := db.Exec("UPDATE cart_detail SET quantity = ? WHERE cartId=? AND itemId=?",
				quantity,
				cart.ID,
				itemId,
			)
			cartDetail.Item.ID = itemId
			cartDetail.Quantity = quantity
			cart.CartDetail = append(cart.CartDetail, cartDetail)
			if errQuery == nil {
				sendSuccessResponse(w, "Update Cart Berhasil", cart)
			} else {
				log.Println(errQuery)
				sendErrorResponse(w, "Update Cart Gagal")
			}
		}
	}
}

// fungsi untuk menghapus item dari cart
func DeleteItemFromCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Failed")
		return
	}
	vars := mux.Vars(r)
	itemId := vars["item_id"]
	_, UserID, _, _ := validateTokenFromCookies(r)
	//jika user belum login, maka akan direturn unauthorized response
	if UserID == -1 {
		sendUnauthorizedResponse(w)
		return
	}
	query := "SELECT cartId FROM cart WHERE userId ='" + strconv.Itoa(UserID) + "'"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var cart model.Cart
	var cartDetail model.CartDetail
	for rows.Next() {
		if err := rows.Scan(&cart.ID); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		} else {
			_, errQuery := db.Exec("DELETE FROM cart_detail WHERE cartId=? AND itemId=?",
				cart.ID,
				itemId,
			)
			cartDetail.Item.ID, _ = strconv.Atoi(itemId)
			cart.CartDetail = append(cart.CartDetail, cartDetail)
			if errQuery == nil {
				sendSuccessResponse(w, "Delete Item dari Cart Berhasil", cart)
			} else {
				log.Println(errQuery)
				sendErrorResponse(w, "Delete Item dari Cart Gagal")
			}
		}
	}
}
