package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	shopaddress := r.Form.Get("shop_address")
	shoptelephone := r.Form.Get("shop_telephone")
	shopemail := r.Form.Get("shop_email")
	shopstatus := r.Form.Get("shop_status")
	query := "UPDATE shop SET "
	if shopname != "" {
		query += "shopname = '" + shopname + "'"
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
		query += "shopcategory = '" + shopcategory + "'"
	}
	if shopaddress != "" {
		if strings.Contains(query, "shopcategory") || strings.Contains(query, "shopreputation") || strings.Contains(query, "shopname") {
			query += ", "
		}
		query += "shopaddress = '" + shopaddress + "'"
	}
	if shoptelephone != "" {
		if strings.Contains(query, "shoptelephone") || strings.Contains(query, "shopcategory") || strings.Contains(query, "shopreputation") || strings.Contains(query, "shopname") {
			query += ", "
		}
		query += "shoptelephone = '" + shoptelephone + "'"
	}
	if shopemail != "" {
		if strings.Contains(query, "shopemail") || strings.Contains(query, "shoptelephone") || strings.Contains(query, "shopcategory") || strings.Contains(query, "shopreputation") || strings.Contains(query, "shopname") {
			query += ", "
		}
		query += "shopemail = '" + shopemail + "'"
	}
	if shopstatus != "" {
		if strings.Contains(query, "shopstatus") || strings.Contains(query, "shopemail") || strings.Contains(query, "shoptelephone") || strings.Contains(query, "shopcategory") || strings.Contains(query, "shopreputation") || strings.Contains(query, "shopname") {
			query += ", "
		}
		query += "shopstatus = " + shopstatus
	}
	query += " WHERE shopid = " + shopid
	_, errQuery := db.Exec(query)
	fmt.Println(query)
	if errQuery != nil {
		log.Println(errQuery)
		sendErrorResponse(w, "Failed to update shop profile")
		return
	} else {
		sendSuccessResponse(w, "Shop profile updated", nil)
	}
}

func GetShopProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	shopid := r.URL.Query().Get("shop_id")
	shopname := r.URL.Query().Get("shop_name")
	shopcategory := r.URL.Query().Get("shop_category")
	shopreputation := r.URL.Query().Get("shop_reputation")

	fmt.Println(shopid, shopname, shopcategory, shopreputation)

	query := "SELECT * FROM shop "

	if shopid != "" {
		query += "WHERE shopid = " + shopid + " "
	}
	if shopname != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " shopName LIKE '%" + shopname + "%' "
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
func RegisterShop(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	currentID := getUserIdFromCookie(r)
	shopname := r.Form.Get("shop_name")
	shopcategory := r.Form.Get("shop_category")
	shopaddress := r.Form.Get("shop_address")
	shoptelephone := r.Form.Get("shop_telephone")
	shopemail := r.Form.Get("shop_email")
	inputpassword := r.Form.Get("password")
	password := GetUserPassword(currentID)

	if password != inputpassword {
		sendErrorResponse(w, "Password does not match")
	} else {
		query := "INSERT INTO shop (shopName,shopReputation,shopCategory,shopAddress,shopTelephone,shopEmail,shopStatus) "
		query += "VALUES (?,?,?,?,?,?,?)"
		res, errQuery := db.Exec(query, shopname, 0, shopcategory, shopaddress, shoptelephone, shopemail, 0)
		if errQuery != nil {
			log.Println(errQuery)
			sendErrorResponse(w, "Failed to register shop")
			return
		} else {
			id, _ := res.LastInsertId()
			shopid := int(id)
			//Ambil dari cookie
			userid := 5
			query = "INSERT INTO shop_admin VALUES(?,?)"
			_, errQuery2 := db.Exec(query, shopid, userid)
			if errQuery2 != nil {
				log.Println(errQuery2)
				sendErrorResponse(w, "Failed to register shop")
			} else {
				var user model.User
				errQuery3 := db.QueryRow("SELECT u.userid, u.name, u.email FROM users u INNER JOIN shop_admin sa ON sa.userId = u.userid INNER JOIN shop s ON s.shopId=sa.shopId WHERE u.userid = ?", userid).Scan(&user.ID, &user.Name, &user.Email)
				if errQuery3 != nil {
					log.Println(errQuery3)
					sendErrorResponse(w, "Failed to register shop")
				} else {
					sendSuccessResponse(w, "Successfully registered shop", nil)
					sendMailRegisShop(user, shopemail)
				}
			}

		}
	}
}
func InsertShopAdmin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	shopid := r.Form.Get("shop_id")
	email := r.Form.Get("email")

	userid := 0
	emailresult := ""
	usertype := 0
	query := "SELECT userid,email,usertype FROM users WHERE email = '" + email + "'"
	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		sendErrorResponse(w, "User with this email doesn't exists!")
		return
	} else {
		for rows.Next() {
			if err := rows.Scan(&userid, &emailresult, &usertype); err != nil {
				log.Println(err)
				sendErrorResponse(w, "Something went wrong, please try again")
				return
			} else {
				if usertype != 2 {
					sendErrorResponse(w, "User with this email is not a seller!")
				} else {
					query2 := "INSERT INTO shop_admin(shopId, userId) VALUES (?,?)"
					_, errQuery2 := db.Exec(query2, shopid, userid)
					if errQuery2 != nil {
						log.Println(errQuery2)
						sendErrorResponse(w, "Failed to add shop admin")
					} else {
						var user model.User
						errQuery := db.QueryRow("SELECT u.userid,u.name,u.email FROM shop_admin sa INNER JOIN users u ON sa.userId=u.userid WHERE sa.shopId=? AND u.userid=?", shopid, userid).Scan(&user.ID, &user.Name, &user.Email)
						if errQuery != nil {
							log.Println(errQuery)
							sendErrorResponse(w, "Failed to register shop")
						} else {
							sendSuccessResponse(w, "Successfully added shop admin", nil)
							sendMailInsertAdmin(user)
						}
					}
				}
			}
		}
	}
}
func GetUserShop(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	currentID := getUserIdFromCookie(r)

	query := "SELECT shop.shopid,shopname,shopreputation,shopcategory,shopaddress,shoptelephone,shopemail FROM shop "
	query += "INNER JOIN shop_admin ON shop.shopid = shop_admin.shopid "
	query += "WHERE shop_admin.userid = " + strconv.Itoa(currentID)
	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	var shop model.Shop
	var shopList []model.Shop
	for rows.Next() {
		if err := rows.Scan(&shop.ID, &shop.Name, &shop.Reputation, &shop.Category, &shop.Address, &shop.TelephoneNo, &shop.Email); err != nil {
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
