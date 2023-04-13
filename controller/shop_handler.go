package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"encoding/json"
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

// get all transaction , filter : transactionID, userID
func GetAllTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	//user jika blm ada cookie
	vars := mux.Vars(r)
	UserID := vars["user_id"]

	//user jika sudah ada cookie
	//_, id, _, _ := validateTokenFromCookies(r)
	//Read From Query Param
	transactionId := r.URL.Query()["transactionId"]
	address := r.URL.Query()["address"]
	date := r.URL.Query()["date"]
	delivery := r.URL.Query()["delivery"]
	progress := r.URL.Query()["progress"]
	paymentType := r.URL.Query()["paymentType"]

	query := "SELECT a.transactionId,a.userId,a.address,a.date,a.delivery,a.progress,a.paymentType,b.transactionId,b.itemId,b.quantity FROM `transaction` a INNER JOIN transaction_detail b ON a.transactionId = b.transactionId"
	if transactionId != nil {
		query += " WHERE transactionId='" + transactionId[0] + "'"
	}
	if address != nil {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		fmt.Println(address[0])
		query += " address LIKE '%" + address[0] + "%' "
	}
	if date != nil {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		fmt.Println(date[0])
		query += " date >= '" + date[0] + "' "
	}
	if delivery != nil {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		fmt.Println(address[0])
		query += " delivery ='" + delivery[0] + "' "
	}
	if progress != nil {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		fmt.Println(progress[0])
		query += " progress='" + progress[0] + "' "
	}
	if paymentType != nil {
		if strings.Contains(query, "WHERE") {
			query += "AND"
		} else {
			query += "WHERE"
		}
		fmt.Println(paymentType[0])
		query += " paymentType='" + paymentType[0] + "' "
	}
	query += " WHERE userId = " + UserID
	fmt.Print(query)
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var transaction model.Transaction
	var transactions []model.Transaction
	var transactionDetail model.TransactionDetail
	var transactionDetails []model.TransactionDetail
	var item model.Item
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.Address, &transaction.Date, &transaction.Delivery, &transaction.Progress, &transaction.PaymentType, &transactionDetail.Quantity, &item.ID, &item.ShopID, &item.Name, &item.Desc, &item.Category, &item.Price, &item.Stock); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		} else {
			//belum ambil isi transaction details dari transactionId
			transactionDetail.Item = item
			transactionDetail.IdTransaction = transaction.ID
			transactionDetails = append(transactionDetails, transactionDetail)
			transactions = append(transactions, transaction)
		}
	}
	for i := range transactions {
		for j := range transactionDetails {
			if transactions[i].ID == transactionDetails[j].IdTransaction {
				transactions[i].TransactionDetail = append(transactions[i].TransactionDetail, transactionDetails[j])
			}
		}
	}
	sendSuccessResponse(w, "Success", transactions)
}

// insert item ke transaction... asumsi insert itemnya itu item yang tidak ada di transaction, kalau itemnya ada, berarti pakai update
// masih memakai user dummy, kalau sudah ada cookie, maka akan diganti cookie
func InsertItemToTransaction(w http.ResponseWriter, r *http.Request) {
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

	//user jika blm ada cookie
	vars := mux.Vars(r)
	UserID := vars["user_id"]

	//user jika sudah ada cookie
	//_, id, _, _ := validateTokenFromCookies(r)
	var transaction model.UpdateTransaction

	_, errQuery := db.Exec("INSERT INTO transaction(transactionId,userId)values(?,?)",
		transaction.TransactionID,
		UserID,
	)
	var response model.GenericResponse
	if errQuery != nil {
		response.Status = 400
		response.Message = "Gagal Insert Cart"
		return
	} else {
		_, errQuery := db.Exec("INSERT INTO transaction_detail(transactionId,itemId,quantity)values(?,?,?)",
			transaction.TransactionID,
			itemId,
			quantity,
		)
		transaction.ItemID = itemId
		transaction.Quantity = quantity
		if errQuery != nil {
			fmt.Println(errQuery)
			response.Status = 400
			response.Message = "Insert Item ke Cart Gagal"
		} else {
			response.Status = 200
			response.Message = "Insert Item ke Cart Berhasil"
			response.Data = transaction
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	vars := mux.Vars(r)
	transId := vars["transaction_id"]
	progress := r.Form.Get("progress")
	sqlStatement := `Update transaction
	SET progress = ?
	where ID =?`

	_, errQuery := db.Exec(sqlStatement, progress, transId)

	if errQuery != nil {
		sendErrorResponse(w, "Failed to update transaction")
		return
	} else {
		sendSuccessResponse(w, "Progress updated", nil)
	}
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
