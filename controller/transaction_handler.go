package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"database/sql"
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
	UserID := getUserIdFromCookie(r)

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
		query2 := "SELECT userid from shop_admin WHERE shopId =? AND userId=?"
		row2 := db.QueryRow(query2, shopId[0], UserID)
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
		//dDsini query transaksi apa saja yang terjadi di toko berdasarkan shopId
		query += " INNER JOIN transaction_detail td ON t.`transactionId` = td.`transactionId` INNER JOIN item i ON td.itemId = i.itemId WHERE i.shopId ='" + shopId[0] + "'"
		if transactionId != nil {
			query += " AND transactionId='" + transactionId[0] + "'"
		}
		query += " GROUP BY t.`transactionId`;"
	}
	fmt.Println(query)
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
	//Kalau query transaksi kosong, ini berarti user yg mengakses bukan pemilik transaksi ini
	//dan beri unauthorized access
	if len(transactions) == 0 {
		if transactionId != nil {
			if transactionId[0] != "" {
				sendUnauthorizedResponse(w)
				return
			}
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
	sendSuccessResponse(w, "Success", transactions)
}

// insert item ke transaction... asumsi insert itemnya itu item yang tidak ada di transaction, kalau itemnya ada, berarti pakai update
func InsertItemToTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	//baca userID dari cookie
	UserID := getUserIdFromCookie(r)
	//baca dari request body
	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Failed")
		return
	}
	address := r.Form.Get("address")
	delivery := r.Form.Get("delivery")
	paymentType := r.Form.Get("payment_type")
	itemIds := r.Form["itemid"]
	quantities := r.Form["quantity"]
	if len(itemIds) != len(quantities) {
		sendErrorResponse(w, "Number of itemId does not match quantity")
		return
	}
	for i := 0; i < len(itemIds); i++ {
		//cek stok terlebih dahulu dengan query stok item yang akan diedit
		query := "SELECT itemStock FROM item WHERE itemId='" + itemIds[i] + "'"
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
		//1 produk transaksi tidak cukup stoknya, transaksi tidak dapar dilakukan
		stock, _ := strconv.Atoi(quantities[i])
		if dbstock < stock {
			sendErrorResponse(w, "Product stock is not enough")
			return
		}
	}
	//Insert transaksi baru
	res, errQuery := db.Exec("INSERT INTO transaction(userId,address,delivery,paymentType)values(?,?,?,?)",
		UserID, address, delivery, paymentType,
	)
	if errQuery != nil {
		sendErrorResponse(w, "Failed while inserting transaction")
		return
	}
	//ID transaction yang baru diinsert
	insertedTransId, _ := res.LastInsertId()
	for i := 0; i < len(itemIds); i++ {
		quantity, err := strconv.Atoi(quantities[i])
		if err != nil {
			sendErrorResponse(w, "Invalid quantity")
			return
		}
		//insert ke detail transaksi
		_, errQuery := db.Exec("INSERT INTO transaction_detail(transactionId,itemId,quantity)values(?,?,?)",
			insertedTransId,
			itemIds[i],
			quantity,
		)
		if errQuery != nil {
			fmt.Println(errQuery)
			sendErrorResponse(w, "Failed to record transaction")
		} else {
			boolStock := reduceStock(itemIds[i], quantity)
			if !boolStock {
				sendErrorResponse(w, "Failed to Update Stock")
			}
		}
	}
	//Update reputasi toko
	UpdateReputation(CheckItemShop(itemIds[0]))
	sendSuccessResponse(w, "Transaction recorded", nil)
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
	sqlStatement := "UPDATE transaction SET progress = ? WHERE transactionID =?"

	_, errQuery := db.Exec(sqlStatement, progress, transId)

	if errQuery != nil {
		sendErrorResponse(w, "Failed to update transaction")
		return
	} else {
		sendSuccessResponse(w, "Transaction progress updated", nil)
	}
}

func reduceStock(itemId string, quantity int) bool {
	//Reduce stock produk setelah transaksi sukses
	db := connect()
	defer db.Close()
	query := "SELECT itemStock FROM item WHERE itemId =?"
	row := db.QueryRow(query, itemId)
	var stock int
	if err := row.Scan(&stock); err != nil {
		return false
	}
	stock = stock - quantity
	_, errQuery := db.Exec("UPDATE item SET itemStock = ? WHERE itemId=?",
		stock,
		itemId,
	)
	if errQuery != nil {
		return false
	} else {
		return true
	}

}
