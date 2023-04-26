package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"database/sql"
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

// insert item ke cart... asumsi insert itemnya itu item yang tidak ada di cart,
// kalau itemnya ada, berarti pakai update
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
	var cart model.Cart
	var cartDetail model.CartDetail
	//dapatkan cartID dari database menggunakan fungsi
	cart.ID = getCartIDFromDatabase(w, UserID)
	// kalau cartID nya itu -1, berarti ada yang eror sehingga langsung return saja
	if cart.ID == -1 {
		return
	}
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
	//cek quantity produk, jika lebih kecil atau sama dengan nol, maka akan direturn response error
	if quantity <= 0 {
		sendErrorResponse(w, "Quantity tidak boleh lebih kecil atau sama dengan nol")
		return
	}
	_, UserID, _, _ := validateTokenFromCookies(r)
	//jika user belum login, maka akan direturn unauthorized response
	if UserID == -1 {
		sendUnauthorizedResponse(w)
		return
	}
	//cek stok terlebih dahulu dengan query stok item yang akan diedit
	query := "SELECT itemStock FROM item WHERE itemId='" + strconv.Itoa(itemId) + "'"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var item model.Item
	for rows.Next() {
		if err := rows.Scan(&item.Stock); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		}
	}
	//jika stok < quantity, maka akan direturn pesan eror
	if item.Stock < quantity {
		sendErrorResponse(w, "Stok kurang dari Quantity")
		return
	}
	var cart model.Cart
	var cartDetail model.CartDetail
	//dapatkan cartID dari database menggunakan fungsi
	cart.ID = getCartIDFromDatabase(w, UserID)
	// kalau cartID nya itu -1, berarti ada yang eror sehingga langsung return saja
	if cart.ID == -1 {
		return
	}
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

	var cart model.Cart
	var cartDetail model.CartDetail
	//dapatkan cartID dari database menggunakan fungsi
	cart.ID = getCartIDFromDatabase(w, UserID)
	// kalau cartID nya itu -1, berarti ada yang eror sehingga langsung return saja
	if cart.ID == -1 {
		return
	}
	//hapus dari tabel cart_detail yang memiliki cartId ... dan itemId ...
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
func getCartIDFromDatabase(w http.ResponseWriter, userID int) int {
	db := connect()
	defer db.Close()

	query := "SELECT cartId FROM cart WHERE userId =?"
	row := db.QueryRow(query, userID)
	var cartId int
	switch err := row.Scan(&cartId); err {
	case sql.ErrNoRows:
		sendErrorResponse(w, "User tidak mempunyai keranjang")
		return -1
	case nil:
		return userID
	default:
		sendErrorResponse(w, "Terjadi kesalahan saat mengecek cartId")
		return -1
	}
}
