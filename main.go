package main

import (
	"fmt"
	"log"
	"net/http"
	"week6/controller"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	/*List usertype
	0 = admin
	1 = staff
	2 = customer
	*/

	//Get
	router.HandleFunc("/users", controller.GetAllUsers).Methods("GET")
	router.HandleFunc("/products", controller.GetAllProducts).Methods("GET")
	router.HandleFunc("/transactions", controller.GetAllTransactions).Methods("GET")

	//Insert
	router.HandleFunc("/users", controller.Authenticate(controller.InsertUser, 0)).Methods("POST")
	router.HandleFunc("/products", controller.InsertProduct).Methods("POST")
	router.HandleFunc("/transactions", controller.InsertTransaction).Methods("POST")

	//Update
	router.HandleFunc("/users/{user_id}", controller.UpdateUser).Methods("PUT")
	router.HandleFunc("/products/{product_id}", controller.UpdateProduct).Methods("PUT")
	router.HandleFunc("/transactions/{transaction_id}", controller.UpdateTransaction).Methods("PUT")

	//Delete
	router.HandleFunc("/users/{user_id}", controller.Authenticate(controller.DeleteUser, 0)).Methods("DELETE")
	router.HandleFunc("/products/{product_id}", controller.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/transactions/{transaction_id}", controller.DeleteTransaction).Methods("DELETE")

	//Detail user transaction
	router.HandleFunc("/usertransactions", controller.GetDetailedTransaction).Methods("GET")
	//Login
	router.HandleFunc("/login", controller.LoginUser).Methods("POST")
	//Logout
	router.HandleFunc("/logout", controller.LogoutUser).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"localhost:8492"},
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)

	http.Handle("/", router)
	fmt.Println("Connected to port 8492")
	log.Println("Connected to port 8492")
	log.Fatal(http.ListenAndServe(":8492", handler))
}
