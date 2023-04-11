package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"week6/model"

	"github.com/gorilla/mux"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	//fmt.Println("Masuk ke GetAllTransactions")

	query := "SELECT * FROM transactions"

	id := r.URL.Query()["id"]
	if id != nil {
		query += " WHERE id = '" + id[0] + "'"
	}

	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		return
	}

	var transaction model.Transaction
	var transactions []model.Transaction

	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity); err != nil {
			log.Println(err)
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}
	var response model.TransactionResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = transactions
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	userID, _ := strconv.Atoi(r.Form.Get("userid"))
	productID, _ := strconv.Atoi(r.Form.Get("productid"))
	quantity, _ := strconv.Atoi(r.Form.Get("qty"))

	_, errQueryInsert := db.Exec("INSERT INTO products(id,name,price) VALUES (?,?,?)", productID, "temp-product", 0)

	res, errQuery := db.Exec("INSERT INTO transactions(userid,productid,quantity) VALUES (?,?,?)", userID, productID, quantity)
	id, _ := res.LastInsertId()

	var response model.TransactionResponse
	if errQuery == nil && errQueryInsert == nil {
		var transactions []model.Transaction
		response.Status = 200
		response.Message = "Insert Success"
		id := int(id)
		transactions = append(transactions, model.Transaction{ID: id, UserID: userID, ProductID: productID, Quantity: quantity})
		response.Data = transactions
	} else {
		response.Status = 400
		response.Message = "Insert Failed"
	}
	w.Header().Set("Content=Type", "application/json")
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
	transactionID := vars["transaction_id"]
	userID, _ := strconv.Atoi(r.Form.Get("userid"))
	productID, _ := strconv.Atoi(r.Form.Get("productid"))
	quantity, _ := strconv.Atoi(r.Form.Get("qty"))
	_, errQuery := db.Exec("UPDATE transactions SET userid = ?, productid = ?, quantity = ? WHERE id = ?",
		userID, productID, quantity, transactionID)

	var response model.TransactionResponse
	if errQuery == nil {
		var transactions []model.Transaction
		response.Status = 200
		response.Message = "Update Success"
		id, _ := strconv.Atoi(transactionID)
		transactions = append(transactions, model.Transaction{ID: id, UserID: userID, ProductID: productID, Quantity: quantity})
		response.Data = transactions
	} else {
		response.Status = 400
		response.Message = "Update Failed"
	}
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	vars := mux.Vars(r)
	transactionID := vars["transaction_id"]
	_, errQuery := db.Exec("DELETE FROM transactions WHERE id = ?", transactionID)

	var response model.TransactionResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
	} else {
		response.Status = 400
		response.Message = "Delete Failed"
	}
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDetailedTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT transactions.id,users.id,users.name,users.age,users.address,users.email,"
	query += "products.id,products.name,products.price,transactions.quantity FROM transactions "
	query += "INNER JOIN users ON users.id = transactions.userid INNER JOIN products ON products.id = transactions.productid"
	id := r.URL.Query()["userid"]
	if id != nil {
		query += " WHERE users.id = '" + id[0] + "'"
	}

	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		return
	}

	var transaction model.DetailedTransaction
	var transactions []model.DetailedTransaction

	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.User.ID, &transaction.User.Name, &transaction.User.Age, &transaction.User.Address,
			&transaction.User.Email, &transaction.Product.ID, &transaction.Product.Name, &transaction.Product.Price,
			&transaction.Quantity); err != nil {
			log.Println(err)
			return
		} else {
			transaction.User.Password = "********"
			transactions = append(transactions, transaction)
		}
	}
	var response model.DetailedTransactionResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = transactions
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
