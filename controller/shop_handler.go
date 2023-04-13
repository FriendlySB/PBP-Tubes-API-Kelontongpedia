package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func UpdateShopProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	vars := mux.Vars(r)
	shopid := vars["shop_id"]
	shopname := r.Form.Get("shop_name")
	shopreputation := r.Form.Get("shop_reputation")
	shopcategory := r.Form.Get("shop_category")
	shopadress := r.Form.Get("shop_adress")
	shoptelephone := r.Form.Get("shop_telephone")
	shopemail := r.Form.Get("shop_email")
	shopstatus := r.Form.Get("shop_status")
	query := "UPDATE shop SET "
	if shopname != "" {
		query += "shopname = " + shopname
	}
	if shopreputation != "" {
		if strings.Contains(query, "shopname") {
			query += ", "
		}
		query += "shopreputation = " + shopreputation
	}
	if shopcategory != "" {
		if strings.Contains(query, "shopreputation") || strings.Contains(query, "shopname") {
			query += ", "
		}
		query += "shopcategory = " + shopcategory
	}
	if shopadress != "" {
		if strings.Contains(query, "shopcategory") || strings.Contains(query, "shopreputation") || strings.Contains(query, "shopname") {
			query += ", "
		}
		query += "shopadress = " + shopadress
	}
	if shoptelephone != "" {
		if strings.Contains(query, "shoptelephone") || strings.Contains(query, "shopcategory") || strings.Contains(query, "shopreputation") || strings.Contains(query, "shopname") {
			query += ", "
		}
		query += "shoptelephone = " + shoptelephone
	}
	if shopemail != "" {
		if strings.Contains(query, "shopemail") || strings.Contains(query, "shoptelephone") || strings.Contains(query, "shopcategory") || strings.Contains(query, "shopreputation") || strings.Contains(query, "shopname") {
			query += ", "
		}
		query += "shopemail = " + shopemail
	}
	if shopstatus != "" {
		if strings.Contains(query, "shopstatus") || strings.Contains(query, "shopemail") || strings.Contains(query, "shoptelephone") || strings.Contains(query, "shopcategory") || strings.Contains(query, "shopreputation") || strings.Contains(query, "shopname") {
			query += ", "
		}
		query += "shopstatus = " + shopstatus
	}
	query += " WHERE shopid = " + shopid
	// _, errQuery := db.Exec(query)
	fmt.Println(query)
	// if errQuery != nil {
	// 	sendErrorResponse(w, "Failed to update transaction")
	// 	return
	// } else {
	// 	sendSuccessResponse(w, "Progress updated", nil)
	// }
}

func GetShopProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	shopname := r.URL.Query().Get("shop_name")
	shopcategory := r.URL.Query().Get("shop_category")
	shopreputation := r.URL.Query().Get("shop_reputation")

	query := "SELECT * FROM shop "

	if shopname != "" {
		query += "WHERE shopName LIKE '%" + shopname + "%' "
	}
	if shopcategory != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " shopCategory = '" + shopcategory + "' "
	}
	if shopreputation != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " shopReputation >= '" + shopreputation + "'"
	}
	fmt.Println(query)
	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	var shop model.Shop
	var shopList []model.Shop
	for rows.Next() {
		if err := rows.Scan(&shop.ID, &shop.Name, &shop.Reputation, &shop.Category, &shop.Address, &shop.TelephoneNo, &shop.Email, &shop.ShopBanStatus); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		} else {
			if !shop.ShopBanStatus {
				shopList = append(shopList, shop)
			}

		}
	}
	sendSuccessResponse(w, "Success", shopList)
}
