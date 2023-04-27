package controller

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func connect() *sql.DB {
	_ = godotenv.Load()
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")
	dbaddress := fmt.Sprintf("root:@tcp(127.0.0.1:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_port, db_name)
	db, err := sql.Open("mysql", dbaddress)
	if err != nil {
		log.Fatal()
	}
	return db
}
