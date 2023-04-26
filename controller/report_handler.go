package controller

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
)

func SetMonthlyReportScheduler() {
	//Setiap tanggal 1 di tiap bulan jam 8 GMT+7, kirim report bulanan
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Month(1).At("01:00").Do(func() {
		SendReport()
	})
	scheduler.StartAsync()
}

func SendReport() {
	listid, listemail := GetShopData()

	for i := 0; i < len(listid); i++ {
		//Dapatkan jumlah transaksi total, jumlah produk terjual, dan income dalam 1 bulan terakhir
		transactionCount, productSold, income := calculateTransaction(listid[i])
		//Kirim data tersebut ke email toko
		//Untuk sementara belum
		fmt.Println(transactionCount, productSold, income, listemail[i])
	}
}
func calculateTransaction(shopid int) (int, int, int) {
	db := connect()
	defer db.Close()

	curTime := time.Now().AddDate(-30, 0, 0)
	date := fmt.Sprint(curTime.Format("1945-08-15"))

	transactionCount := 0
	productSold := 0
	income := 0

	query := "SELECT SUM(transaction_detail.quantity),SUM(item.itemPrice*transaction_detail.quantity) FROM transaction "
	query += "INNER JOIN transaction_detail ON transaction.transactionId = transaction_detail.transactionId "
	query += "INNER JOIN item ON item.itemId = transaction_detail.itemId "
	query += "WHERE item.shopId = " + strconv.Itoa(shopid) + " AND date >= '" + date + "' "
	query += "GROUP BY transaction.transactionId "
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	} else {
		tempSold := 0
		tempIncome := 0
		for rows.Next() {
			if err := rows.Scan(&tempSold, &tempIncome); err != nil {
				log.Print(err)
				return 0, 0, 0
			} else {
				transactionCount += 1
				productSold += tempSold
				income += tempIncome
			}
		}
	}
	return transactionCount, productSold, income
}
func GetShopData() ([]int, []string) {
	db := connect()
	defer db.Close()

	var listid []int
	var listemail []string

	query := "SELECT shopid,shopemail FROM shop WHERE shopstatus != 1"

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	} else {
		var id int
		var email string
		for rows.Next() {
			if err := rows.Scan(&id, &email); err != nil {
				log.Print(err)
				return nil, nil
			} else {
				listid = append(listid, id)
				listemail = append(listemail, email)
			}
		}
	}
	return listid, listemail
}
func GetShopAdminEmails(shopid int) []string {
	db := connect()
	defer db.Close()

	var listemail []string

	query := "SELECT email FROM users "
	query += "INNER JOIN shop_admin ON users.userid = shop_admin.userid "
	query += "WHERE shopid = " + strconv.Itoa(shopid)
	fmt.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	} else {
		var email string
		for rows.Next() {
			if err := rows.Scan(&email); err != nil {
				log.Print(err)
				return nil
			} else {
				listemail = append(listemail, email)
			}
		}
	}
	return listemail
}
