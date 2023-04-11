package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"week6/model"

	"github.com/gorilla/mux"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	//fmt.Println("Masuk ke GetAllProducts")

	query := "SELECT * FROM products"

	name := r.URL.Query()["name"]
	if name != nil {
		query += " WHERE name = '" + name[0] + "'"
	}

	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		return
	}

	var product model.Product
	var products []model.Product

	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		} else {
			products = append(products, product)
		}
	}
	var response model.ProductResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = products
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	res, errQuery := db.Exec("INSERT INTO products(name,price) VALUES (?,?)", name, price)
	id, _ := res.LastInsertId()

	var response model.ProductResponse
	if errQuery == nil {
		var products []model.Product
		response.Status = 200
		response.Message = "Insert Success"
		id := int(id)
		products = append(products, model.Product{ID: id, Name: name, Price: price})
		response.Data = products
	} else {
		response.Status = 400
		response.Message = "Insert Failed"
	}
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	vars := mux.Vars(r)
	productID := vars["product_id"]
	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))
	_, errQuery := db.Exec("UPDATE products SET name = ?, price = ? WHERE id = ?", name, price, productID)

	var response model.ProductResponse
	if errQuery == nil {
		var products []model.Product
		var product model.Product
		response.Status = 200
		response.Message = "Update Success"
		id, _ := strconv.Atoi(productID)
		product.Name = name
		product.Price = price
		products = append(products, model.Product{ID: id, Name: name, Price: price})
		response.Data = products
	} else {
		response.Status = 400
		response.Message = "Update Failed"
	}
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	vars := mux.Vars(r)
	productID := vars["product_id"]

	_, errQueryDelete := db.Exec("DELETE FROM transactions WHERE productid = ?", productID)

	_, errQuery := db.Exec("DELETE FROM products WHERE id = ?", productID)

	var response model.ProductResponse
	if errQuery == nil && errQueryDelete == nil {
		response.Status = 200
		response.Message = "Success"
	} else {
		response.Status = 400
		response.Message = "Delete Failed"
	}
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
