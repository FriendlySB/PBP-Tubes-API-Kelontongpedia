package controller

import (
	"fmt"
	"net/http"
	"strings"
)

func UpdateShopProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	shopid := r.Form.Get("shop_id")
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
		if strings.Contains(query, "shopreputation") {
			query += ", "
		}
		query += "shopcategory = " + shopcategory
	}
	if shopadress != "" {
		if strings.Contains(query, "shopcategory") {
			query += ", "
		}
		query += "shopadress = " + shopadress
	}
	if shoptelephone != "" {
		if strings.Contains(query, "shoptelephone") {
			query += ", "
		}
		query += "shoptelephone = " + shoptelephone
	}
	if shopemail != "" {
		if strings.Contains(query, "shopemail") {
			query += ", "
		}
		query += "shopemail = " + shopemail
	}
	if shopstatus != "" {
		if strings.Contains(query, "shopstatus") {
			query += ", "
		}
		query += "shopstatus = " + shopstatus
	}
	query += " WHERE shopid = " + shopid
	//_, errQuery := db.Exec(query)
	fmt.Println(query)
}
