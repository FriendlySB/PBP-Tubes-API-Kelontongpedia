package main

import (
	"PBP-Tubes-API-Tokopedia/controller"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	//List End Points
	/*
		//Untuk database
		//Usertype
			0 = admin
			1 = pembeli
			2 = toko
		//ban
			0 = not banned
			1 = banned
		//Dibagi berdasarkan siapa yang bisa mengaksesnya
		//Umum
		1. Login
		2. Logout
		3. Register
		4. GetAllItem #Dipake di pembeli dan shop juga
		5. GetItemReview #Dipake di pembeli dan shop juga
		6. GetShopProfile #Dipake di pembeli dan shop juga

		//Pembeli
		1. GetUserProfile
		2. UpdateUserProfile
		3. InsertCart
		4. UpdateCart
		5. RemoveCart
		6. ReviewItem
		7. GetTransaction #Dipake di shop juga buat ngebuat daftar penjualan toko

		//Shop
		1. InsertItem
		2. UpdateItem
		3. DeleteItem
		4. UpdateTransaction
		5. UpdateShopProfile

		//Admin
		1. BanUser
		2. BanToko
		3. GetAllUser
	*/

	router.HandleFunc("/login", controller.Login).Methods("POST")
	router.HandleFunc("/logout", controller.Logout).Methods("POST")
	router.HandleFunc("/password", controller.ChangePassword).Methods("PUT")

	router.HandleFunc("/cart", controller.GetCart).Methods("GET")
	router.HandleFunc("/cart", controller.InsertItemToCart).Methods("POST")
	router.HandleFunc("/cart", controller.UpdateCart).Methods("PUT")
	router.HandleFunc("/cart/{item_id}", controller.DeleteItemFromCart).Methods("DELETE")

	router.HandleFunc("/item", controller.GetItem).Methods("GET")
	router.HandleFunc("/item", controller.InsertItem).Methods("POST")
	router.HandleFunc("/item", controller.UpdateItem).Methods("PUT")
	router.HandleFunc("/item/{item_id}", controller.DeleteItem).Methods("Delete")

	router.HandleFunc("/shop", controller.GetShopProfile).Methods("GET")
	router.HandleFunc("/shop", controller.RegisterShop).Methods("POST")
	router.HandleFunc("/shop/{shop_id}", controller.UpdateShopProfile).Methods("PUT")

	router.HandleFunc("/updateTransaction/{transaction_id}", controller.UpdateTransaction).Methods("PUT")

	router.HandleFunc("/profile", controller.GetUserProfile).Methods("GET")
	router.HandleFunc("/updateprofile", controller.UpdateUserProfile).Methods("PUT")

	router.HandleFunc("/review", controller.ReviewItem).Methods("POST")
	router.HandleFunc("/review", controller.GetItemReview).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"localhost:8080"},
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)

	http.Handle("/", router)

	fmt.Println("Connected to port 8181")
	log.Println("Connected to port 8181")

	log.Fatal(http.ListenAndServe(":8181", handler))
}
