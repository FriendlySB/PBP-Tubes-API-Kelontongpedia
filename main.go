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
	router.HandleFunc("/cart", controller.GetCart).Methods("GET")
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
