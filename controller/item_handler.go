package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Select sebuah item dari database
func GetItem(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	itemid := r.URL.Query().Get("item_id")
	itemname := r.URL.Query().Get("item_name")
	itemcategory := r.URL.Query().Get("item_category")
	itemprice := r.URL.Query().Get("item_price")
	shopid := r.URL.Query().Get("shop_id")

	query := "SELECT * FROM item "

	if itemid != "" {
		query += "WHERE itemid = " + itemid + " "
	}
	if itemname != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " itemName LIKE '%" + itemname + "%' "
	}
	if itemcategory != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " itemCategory = '" + itemcategory + "' "
	}
	if itemprice != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " itemPrice <= '" + itemprice + "'"
	}
	if shopid != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " shopid = '" + shopid + "'"
	}
	if strings.Contains(query, "WHERE") {
		query += "AND"
	} else {
		query += "WHERE"
	}
	//Cek agar toko yang diban tidak muncul produknya saat disearch
	query += " shopid NOT IN ("
	query += "SELECT shopid FROM shop WHERE shopstatus = 1)"
	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	var item model.Item
	var itemList []model.Item
	for rows.Next() {
		if err := rows.Scan(&item.ID, &item.ShopID, &item.Name, &item.Desc, &item.Category, &item.Price, &item.Stock); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		} else {
			itemList = append(itemList, item)
		}
	}
	sendSuccessResponse(w, "Success", itemList)
}

// Insert item ke database
func InsertItem(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	shopid := r.Form.Get("shop_id")
	itemname := r.Form.Get("item_name")
	itemdesc := r.Form.Get("item_desc")
	itemcategory := r.Form.Get("item_category")
	itemprice := r.Form.Get("item_price")
	itemstock := r.Form.Get("item_stock")

	UserID := getUserIdFromCookie(r)
	//Cek penjual yang akses ini admin toko yang bersangkutan atau bukan. Jika bukan, unauthorized
	query2 := "SELECT userid from shop_admin WHERE shopId =? AND userId=?"
	row2 := db.QueryRow(query2, shopid, UserID)
	//jika terjadi error saat mengecek, berarti user yg mengakses bukan admin toko ini dan beri
	//unauthorized access
	var temp int
	switch err := row2.Scan(&temp); err {
	case sql.ErrNoRows:
		sendUnauthorizedResponse(w)
		return
	case nil:

	default:
		sendErrorResponse(w, "Error")
		return
	}

	query := "INSERT INTO item (shopid,itemname,itemdesc,itemcategory,itemprice,itemstock) VALUES (?,?,?,?,?,?)"
	_, errQuery := db.Exec(query, shopid, itemname, itemdesc, itemcategory, itemprice, itemstock)

	if errQuery != nil {
		log.Println(err)
		sendErrorResponse(w, "Failed to add new item")
		return
	} else {
		sendSuccessResponse(w, "Successfully added new item", nil)
	}
}
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	vars := mux.Vars(r)

	itemid := vars["item_id"]
	itemname := r.Form.Get("item_name")
	itemdesc := r.Form.Get("item_desc")
	itemcategory := r.Form.Get("item_category")
	itemprice := r.Form.Get("item_price")
	itemstock := r.Form.Get("item_stock")

	UserID := getUserIdFromCookie(r)
	shopid := CheckItemShop(itemid)
	//Cek penjual yang akses ini admin toko yang bersangkutan atau bukan. Jika bukan, unauthorized
	query2 := "SELECT userid from shop_admin WHERE shopId =? AND userId=?"
	row2 := db.QueryRow(query2, shopid, UserID)
	//jika terjadi error saat mengecek, berarti user yg mengakses bukan admin toko ini dan beri
	//unauthorized access
	var temp int
	switch err := row2.Scan(&temp); err {
	case sql.ErrNoRows:
		sendUnauthorizedResponse(w)
		return
	case nil:

	default:
		sendErrorResponse(w, "Error")
		return
	}

	query := "UPDATE item SET "
	if itemname != "" {
		query += "itemname = '" + itemname + "'"
	}
	if itemdesc != "" {
		if query[len(query)-1:] != " " {
			query += ", "
		}
		query += "itemdesc = '" + itemdesc + "'"
	}
	if itemcategory != "" {
		if query[len(query)-1:] != " " {
			query += ", "
		}
		query += "itemcategory = '" + itemcategory + "'"
	}
	if itemprice != "" {
		if query[len(query)-1:] != " " {
			query += ", "
		}
		query += "itemprice = " + itemprice
	}
	if itemstock != "" {
		if query[len(query)-1:] != " " {
			query += ", "
		}
		query += "itemstock = " + itemstock
	}
	query += " WHERE itemid = " + itemid

	_, errQuery := db.Exec(query)

	if errQuery != nil {
		log.Println(errQuery)
		sendErrorResponse(w, "Failed to update item")
		return
	} else {
		sendSuccessResponse(w, "Successfully updated item", nil)
	}
}
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	itemid := vars["item_id"]

	UserID := getUserIdFromCookie(r)
	shopid := CheckItemShop(itemid)
	//Cek penjual yang akses ini admin toko yang bersangkutan atau bukan. Jika bukan, unauthorized
	query2 := "SELECT userid from shop_admin WHERE shopId =? AND userId=?"
	row2 := db.QueryRow(query2, shopid, UserID)
	//jika terjadi error saat mengecek, berarti user yg mengakses bukan admin toko ini dan beri
	//unauthorized access
	var temp int
	switch err := row2.Scan(&temp); err {
	case sql.ErrNoRows:
		sendUnauthorizedResponse(w)
		return
	case nil:

	default:
		sendErrorResponse(w, "Error")
		return
	}

	query := "DELETE FROM item WHERE itemid = " + itemid
	result, errQuery := db.Exec(query)
	if errQuery != nil {
		log.Println(errQuery)
		sendErrorResponse(w, "Failed to delete item")
		return
	} else {
		num, _ := result.RowsAffected()
		if num == 0 {
			sendSuccessResponse(w, "No item deleted", nil)
		} else {
			sendSuccessResponse(w, "Successfully deleted item", nil)
		}
	}
}

func CheckItemShop(itemid string) int {
	db := connect()
	defer db.Close()
	query := "SELECT shopid FROM item WHERE itemid =?"
	row := db.QueryRow(query, itemid)
	var shopid int
	switch err := row.Scan(&shopid); err {
	case sql.ErrNoRows:
		return -1
	case nil:
		return shopid
	default:
		return -1
	}
}

func getShopItem(shopid string) []int {
	db := connect()
	defer db.Close()

	var itemlist []int

	query := "SELECT itemid FROM item WHERE shopid = ?"
	rows, err := db.Query(query, shopid)

	if err != nil {
		log.Println(err)
		return itemlist
	}

	for rows.Next() {
		var itemid int
		if err := rows.Scan(&itemid); err != nil {
			log.Println(err)
			return itemlist
		} else {
			itemlist = append(itemlist, itemid)
		}
	}
	return itemlist
}
