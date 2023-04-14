package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// get all transaction , filter : transactionID, userID
func GetAllTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	//baca userID dari cookie
	_, UserID, _, _ := validateTokenFromCookies(r)
	if UserID == -1 {
		sendUnauthorizedResponse(w)
		return
	}
	//baca dari Query Param
	transactionId := r.URL.Query()["transaction_id"]
	shopId := r.URL.Query()["shop_id"]
	query := "SELECT `transactionId`,`address`,`date`,`delivery`,`progress`,`paymentType` FROM `transaction` "
	if shopId == nil {
		query += " WHERE userId='" + strconv.Itoa(UserID) + "'"
		if transactionId != nil {
			query += " AND transactionId='" + transactionId[0] + "'"
		}
	}
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var transaction model.Transaction
	var transactions []model.Transaction
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.Address, &transaction.Date, &transaction.Delivery, &transaction.Progress, &transaction.PaymentType); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Error result scan")
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}
	for i, v := range transactions {
		query2 := "SELECT a.quantity,b.itemId,b.shopId,b.itemName,b.itemDesc,b.itemCategory,b.itemPrice,b.itemStock FROM transaction_detail a INNER JOIN item b ON a.itemId =b.itemId WHERE a.transactionId = '" + strconv.Itoa(v.ID) + "'"
		if shopId != nil {
			query2 += " AND b.itemId IN (SELECT itemId FROM item where shopId='" + shopId[0] + "')"
		}
		rows2, err2 := db.Query(query2)
		if err2 != nil {
			log.Println(err2)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		}

		var transactionDetail model.TransactionDetail
		var transactionDetails []model.TransactionDetail
		var item model.Item
		for rows2.Next() {
			if err := rows2.Scan(&transactionDetail.Quantity, &item.ID, &item.ShopID, &item.Name, &item.Desc, &item.Category, &item.Price, &item.Stock); err != nil {
				log.Println(err)
				sendErrorResponse(w, "Error result scan")
				return
			} else {
				transactionDetail.Item = item
				transactionDetails = append(transactionDetails, transactionDetail)
			}
		}
		transactions[i].TransactionDetail = transactionDetails
	}
	//biar kalau ada transaksi yang detail transaksinya null dihapus
	var filteredTransactions []model.Transaction
	for _, v := range transactions {
		if v.TransactionDetail != nil {
			filteredTransactions = append(filteredTransactions, v)
		}
	}
	sendSuccessResponse(w, "Get Transaction Berhasil", filteredTransactions)
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
