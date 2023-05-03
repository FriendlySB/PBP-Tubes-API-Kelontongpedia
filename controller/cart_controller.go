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
	UserID := getUserIdFromCookie(r)

	//Query ke database
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

	if len(cartDetails) == 0 {
		sendSuccessResponse(w, "Success", nil)
	} else {
		sendSuccessResponse(w, "Success", cart)
	}

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
	UserID := getUserIdFromCookie(r)
	//dapatkan cartID dari database menggunakan fungsi
	cartid := getCartIDFromDatabase(UserID)
	// kalau cartID nya itu -1, berarti ada yang eror sehingga langsung return saja
	if cartid == -1 {
		sendErrorResponse(w, "Error fetching user cart id")
		return
	}
	qtyInCart := checkItemInCart(cartid, itemId)
	if qtyInCart == 0 {
		_, errQuery := db.Exec("INSERT INTO cart_detail(cartId,itemId,quantity)values(?,?,?)",
			cartid,
			itemId,
			quantity,
		)
		if errQuery == nil {
			sendSuccessResponse(w, "Successfully inserted product to cart", nil)
		} else {
			log.Println(errQuery)
			sendErrorResponse(w, "Failed to insert product to cart")
		}
	} else {
		newquantity := qtyInCart + quantity
		query := "UPDATE cart_detail SET quantity = ? WHERE cartid = ? AND itemid= ?"
		_, errQuery := db.Exec(query, newquantity, cartid, itemId)
		if errQuery == nil {
			sendSuccessResponse(w, "Successfully inserted product to cart", nil)
		} else {
			log.Println(errQuery)
			sendErrorResponse(w, "Failed to insert product to cart")
		}
	}
}

func checkItemInCart(cartid int, itemid int) int {
	db := connect()
	defer db.Close()

	var qty int
	query := "SELECT quantity FROM cart_detail WHERE cartid = ? AND itemid = ?"
	row := db.QueryRow(query, cartid, itemid)
	switch err := row.Scan(&qty); err {
	case sql.ErrNoRows:
		return 0
	case nil:
		return qty
	default:
		return -1
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
		sendErrorResponse(w, "Quantity cannot be less than or equal to zero")
		return
	}
	UserID := getUserIdFromCookie(r)
	//cek stok terlebih dahulu dengan query stok item yang akan diedit
	query := "SELECT itemStock FROM item WHERE itemId='" + strconv.Itoa(itemId) + "'"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var dbstock int
	for rows.Next() {
		if err := rows.Scan(&dbstock); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		}
	}
	//jika stok < quantity, maka akan direturn pesan eror
	if dbstock < quantity {
		sendErrorResponse(w, "Stock is not enough")
		return
	}
	//dapatkan cartID dari database menggunakan fungsi
	cartid := getCartIDFromDatabase(UserID)
	// kalau cartID nya itu -1, berarti ada yang eror sehingga langsung return saja
	if cartid == -1 {
		sendErrorResponse(w, "Error fetching user cart id")
		return
	}
	_, errQuery := db.Exec("UPDATE cart_detail SET quantity = ? WHERE cartId=? AND itemId=?",
		quantity,
		cartid,
		itemId,
	)
	if errQuery == nil {
		sendSuccessResponse(w, "Successfully updated cart product", nil)
	} else {
		log.Println(errQuery)
		sendErrorResponse(w, "Failed to update cart product")
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
	UserID := getUserIdFromCookie(r)

	//dapatkan cartID dari database menggunakan fungsi
	cartid := getCartIDFromDatabase(UserID)
	// kalau cartID nya itu -1, berarti ada yang eror sehingga langsung return saja
	if cartid == -1 {
		sendErrorResponse(w, "Error fetching user cart id")
		return
	}
	//hapus dari tabel cart_detail yang memiliki cartId ... dan itemId ...
	id, _ := strconv.Atoi(itemId)
	delSuccess := removeItemFromCart(cartid, id)

	if delSuccess {
		sendSuccessResponse(w, "Product successfully removed from cart", nil)
	} else {
		sendErrorResponse(w, "Failed to remove product from cart")
	}

}
func getCartIDFromDatabase(userID int) int {
	db := connect()
	defer db.Close()

	query := "SELECT cartId FROM cart WHERE userId =?"
	row := db.QueryRow(query, userID)
	var cartId int
	switch err := row.Scan(&cartId); err {
	case sql.ErrNoRows:
		return -1
	case nil:
		return cartId
	default:
		return -1
	}
}

func removeItemFromCart(cartid int, itemId int) bool {
	db := connect()
	defer db.Close()
	_, errQuery := db.Exec("DELETE FROM cart_detail WHERE cartId=? AND itemId=?",
		cartid,
		itemId,
	)
	if errQuery == nil {
		return true
	} else {
		return false
	}
}

func removeBannedItemFromCart(itemid int) {
	db := connect()
	defer db.Close()
	db.Exec("DELETE FROM cart_detail WHERE itemId=?", itemid)
}
