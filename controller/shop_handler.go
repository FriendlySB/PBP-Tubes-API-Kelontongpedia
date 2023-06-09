package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	CurUserID := getUserIdFromCookie(r)
	vars := mux.Vars(r)
	shopid := vars["shop_id"]
	shopname := r.Form.Get("shop_name")
	shopcategory := r.Form.Get("shop_category")
	shopaddress := r.Form.Get("shop_address")
	shoptelephone := r.Form.Get("shop_telephone")
	shopemail := r.Form.Get("shop_email")

	//Cek apakah penjual admin toko ini. Kalau bukan, unauthorized access
	if !CheckShopAdmin(CurUserID, shopid) {
		sendUnauthorizedResponse(w)
		return
	}

	query := "UPDATE shop SET "
	if shopname != "" {
		query += "shopname = '" + shopname + "'"
	}
	if shopcategory != "" {
		if query[len(query)-1:] != " " {
			query += ", "
		}
		query += "shopcategory = '" + shopcategory + "'"
	}
	if shopaddress != "" {
		if query[len(query)-1:] != " " {
			query += ", "
		}
		query += "shopaddress = '" + shopaddress + "'"
	}
	if shoptelephone != "" {
		if query[len(query)-1:] != " " {
			query += ", "
		}
		query += "shoptelephone = '" + shoptelephone + "'"
	}
	if shopemail != "" {
		if query[len(query)-1:] != " " {
			query += ", "
		}
		query += "shopemail = '" + shopemail + "'"
	}
	query += " WHERE shopid = " + shopid
	_, errQuery := db.Exec(query)
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

	hash := sha256.Sum256([]byte(inputpassword))
	passwordHash := hex.EncodeToString(hash[:])

	if password != passwordHash {
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
			query = "INSERT INTO shop_admin VALUES(?,?)"
			_, errQuery2 := db.Exec(query, shopid, currentID)
			if errQuery2 != nil {
				log.Println(errQuery2)
				sendErrorResponse(w, "Failed to register shop")
			} else {
				var user model.User
				errQuery3 := db.QueryRow("SELECT u.userid, u.name, u.email FROM users u INNER JOIN shop_admin sa ON sa.userId = u.userid INNER JOIN shop s ON s.shopId=sa.shopId WHERE u.userid = ?", currentID).Scan(&user.ID, &user.Name, &user.Email)
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
	CurUserID := getUserIdFromCookie(r)
	shopid := r.Form.Get("shop_id")
	email := r.Form.Get("email")

	//Cek apakah penjual admin toko ini. Kalau bukan, unauthorized access
	if !CheckShopAdmin(CurUserID, shopid) {
		sendUnauthorizedResponse(w)
		return
	}
	userid := 0
	emailresult := ""
	usertype := 0
	query := "SELECT userid,email,usertype FROM users WHERE email = '" + email + "'"
	rows := db.QueryRow(query)
	switch err := rows.Scan(&userid, &emailresult, &usertype); err {
	case sql.ErrNoRows:
		sendErrorResponse(w, "User with this email not found")
		return
	case nil:
		if usertype != 2 {
			sendErrorResponse(w, "User with this email is not a seller!")
			return
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
	default:
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
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
func CheckShopAdmin(userid int, shopid string) bool {
	db := connect()
	defer db.Close()

	query := "SELECT shopid FROM shop_admin WHERE userid=? AND shopid = ?"
	row := db.QueryRow(query, userid, shopid)
	var temp int
	switch err := row.Scan(&temp); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		return false
	}
}
func UpdateReputation(shopid int) {
	db := connect()
	defer db.Close()
	//Update reputasi toko. Setiap transaksi sukses akan menambah 5 poin reputasi
	query := "SELECT shopreputation FROM shop WHERE shopid = ?"
	row := db.QueryRow(query, shopid)
	var reputation int
	if err := row.Scan(&reputation); err != nil {
		return
	}
	reputation = reputation + 5
	db.Exec("UPDATE shop SET shopreputation = ? WHERE shopid=?", reputation, shopid)
}
