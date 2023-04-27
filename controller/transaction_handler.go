package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// get all transaction , filter : transactionID dan  userID , shopID dan userID
func GetAllTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	//baca userID dari cookie
	_, UserID, _, _ := validateTokenFromCookies(r)
	//jika user belum login, maka akan direturn unauthorized response
	if UserID == -1 {
		sendUnauthorizedResponse(w)
		return
	}
	//baca dari Query Param
	transactionId := r.URL.Query()["transaction_id"]
	shopId := r.URL.Query()["shop_id"]
	//query untuk mengambil daftar transaksi yang akan dibaca
	query := "SELECT DISTINCT t.transactionId, t.address, t.date, t.delivery, t.progress, t.paymentType FROM transaction t"
	//jika tidak ada shopId, berarti ini query pengguna, maka akan diambil transaksi pengguna
	//dan jika ada transactonId, transaksi itu saja yang akan diambil
	if shopId == nil {
		query += " WHERE userId='" + strconv.Itoa(UserID) + "'"
		if transactionId != nil {
			query += " AND transactionId='" + transactionId[0] + "'"
		}
	} else {
		//jika ada shopId, maka id user yang mengakses harus dicek, apakah terdaftar di daftar admin toko.
		var check bool
		query2 := "select * from shop_admin WHERE shopId=? AND userId=?"
		row2 := db.QueryRow(query2, shopId[0], UserID)
		err2 := row2.Scan(&check)
		//jika terjadi eror saat mengecek
		if err2 != nil {
			sendErrorResponse(w, "Error ketika mengecek apakah user admin toko")
			return
		}
		//jika user bukan admin toko, maka akan direturn unauthorized response
		if !check {
			sendUnauthorizedResponse(w)
			return
		}
		//disini harusnya query transaksi apa saja yang terjadi di toko berdasarkan shopId, tapi belum beres
		query = " INNER JOIN transaction_detail td ON t.`transactionId` = td.`transactionId` INNER JOIN item i ON td.itemId = i.itemId WHERE i.shopId ='" + shopId[0] + "'"
		if transactionId != nil {
			query += " AND transactionId='" + transactionId[0] + "'"
		}
		query += " GROUP BY t.`transactionId`;"
	}
	//eksekusi query mencari transaksi apa saja yang ada
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	//masukan transaksi tersebut ke array transactions
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
	//query setiap item yang ada di transaksi tersebut
	for i, v := range transactions {
		query2 := "SELECT a.quantity,b.itemId,b.shopId,b.itemName,b.itemDesc,b.itemCategory,b.itemPrice,b.itemStock FROM transaction_detail a INNER JOIN item b ON a.itemId =b.itemId WHERE a.transactionId = '" + strconv.Itoa(v.ID) + "'"
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
	//dikomen dulu karena blm dicek lagi
	//biar kalau ada transaksi yang detail transaksinya null dihapus
	// var filteredTransactions []model.Transaction
	// for _, v := range transactions {
	// 	if v.TransactionDetail != nil {
	// 		filteredTransactions = append(filteredTransactions, v)
	// 	}
	// }
	sendSuccessResponse(w, "Get Transaction Berhasil", transactions)
}

// insert item ke transaction... asumsi insert itemnya itu item yang tidak ada di transaction, kalau itemnya ada, berarti pakai update
func InsertItemToTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	//baca userID dari cookie
	_, UserID, _, _ := validateTokenFromCookies(r)
	//jika user belum login, maka akan direturn unauthorized response
	if UserID == -1 {
		sendUnauthorizedResponse(w)
		return
	}
	//baca dari request body
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed")
		return
	}
	itemIds := r.Form["itemId"]
	quantities := r.Form["quantity"]
	if len(itemIds) != len(quantities) {
		sendErrorResponse(w, "Jumlah ItemId tidak sama dengan Quantity")
		return
	}

	var transaction model.Transaction
	//Insert transaksi baru
	_, errQuery := db.Exec("INSERT INTO transaction(transactionId,userId)values(?,?)",
		transaction.ID,
		UserID,
	)
	if errQuery != nil {
		sendErrorResponse(w, "gagal insert transaksi")
		return
		for i, itemId := range itemIds {
			quantity, err := strconv.Atoi(quantities[i])
			if err != nil {
				sendErrorResponse(w, "Invalid quantity")
				return
			}
			//insert ke detail transaksi
			_, errQuery := db.Exec("INSERT INTO transaction_detail(transactionId,itemId,quantity)values(?,?,?)",
				transaction.ID,
				itemId,
				quantity,
			)
			var transactionDetail model.TransactionDetail
			//broken disini
			transactionDetail.Item.ID = itemId
			transactionDetail.Quantity = quantity
			transaction.TransactionDetail = append(transaction.TransactionDetail, transactionDetail)
			if errQuery != nil {
				fmt.Println(errQuery)
				sendErrorResponse(w, "gagal insert item ke transaksi")
			} else {
				sendSuccessResponse(w, "Insert item ke transaksi berhasil", transaction)
			}
		}
	}
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
