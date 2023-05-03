package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"bytes"
	"fmt"
	"strconv"
	"text/template"

	gm "gopkg.in/gomail.v2"
)

func sendMailRegis(user model.User) {
	mail := gm.NewMessage()

	template := "bin/template/mailRegis.html"

	result, _ := parseTemplate(template, user)

	mail.SetHeader("From", "kelontongpedia2023@gmail.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Notifications")
	mail.SetBody("text/html", result)

	sender := gm.NewDialer("smtp.gmail.com", 25, "kelontongpedia2023@gmail.com", "vdfsiejrvbjrpnyg")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent to: ", user.Email)
	}
}

// func sendMailLogin(user model.User) {
// 	mail := gm.NewMessage()

// 	template := "bin/template/mailLogin.html"

// 	result, _ := parseTemplate(template, user)

// 	mail.SetHeader("From", "kelontongpedia2023@gmail.com")
// 	mail.SetHeader("To", user.Email)
// 	mail.SetHeader("Subject", "Notifications")
// 	mail.SetBody("text/html", result)

// 	sender := gm.NewDialer("smtp.gmail.com", 25, "kelontongpedia2023@gmail.com", "vdfsiejrvbjrpnyg")

// 	if err := sender.DialAndSend(mail); err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println("Email sent to: ", user.Email)
// 	}
// }

func sendMailBanUser(user model.User) {
	mail := gm.NewMessage()

	template := "bin/template/mailBanUser.html"

	result, _ := parseTemplate(template, user)

	mail.SetHeader("From", "kelontongpedia2023@gmail.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Notifications")
	mail.SetBody("text/html", result)

	sender := gm.NewDialer("smtp.gmail.com", 25, "kelontongpedia2023@gmail.com", "xhdfxoyciurbfizb")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent to: ", user.Email)
	}
}

func sendMailRegisShop(user model.User, shop string) {
	// Fungsi untuk mengirim email
	sendEmail := func(to, template string) {
		mail := gm.NewMessage()
		result, _ := parseTemplate(template, user)
		mail.SetHeader("From", "kelontongpedia2023@gmail.com")
		mail.SetHeader("To", to)
		mail.SetHeader("Subject", "Notifications")
		mail.SetBody("text/html", result)
		sender := gm.NewDialer("smtp.gmail.com", 25, "kelontongpedia2023@gmail.com", "xhdfxoyciurbfizb")
		if err := sender.DialAndSend(mail); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Email sent to: ", to)
		}
	}

	template := "bin/template/mailRegisShop.html"
	// Mengirim email pertama ke user
	sendEmail(user.Email, template)

	// Mengirim email kedua ke shop
	sendEmail(shop, template)
}

func sendMailInsertAdmin(user model.User) {
	mail := gm.NewMessage()

	template := "bin/template/mailInsertAdmin.html"

	result, _ := parseTemplate(template, user)

	mail.SetHeader("From", "kelontongpedia2023@gmail.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Notifications")
	mail.SetBody("text/html", result)

	sender := gm.NewDialer("smtp.gmail.com", 25, "kelontongpedia2023@gmail.com", "xhdfxoyciurbfizb")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent to: ", user.Email)
	}
}

func sendMailBanShop(shop model.Shop) {
	mail := gm.NewMessage()

	template := "bin/template/mailBanShop.html"

	result, _ := parseTemplate(template, shop)

	mail.SetHeader("From", "kelontongpedia2023@gmail.com")
	mail.SetHeader("To", shop.Email)
	mail.SetHeader("Subject", "Notifications")
	mail.SetBody("text/html", result)

	sender := gm.NewDialer("smtp.gmail.com", 25, "kelontongpedia2023@gmail.com", "xhdfxoyciurbfizb")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent to: ", shop.Email)
	}
}

func sendMailMonthlyReport(transactionCount int, productSold int, income int, email string) {
	mail := gm.NewMessage()

	mail.SetHeader("From", "kelontongpedia2023@gmail.com")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Monthly Report")
	body := "<h1>Monthly Shop Report</h1>"
	body += "<h2>Hello there</h2>"
	body += "<p>This is your shop's monthly report which includes the following:</p>"
	body += "<ul><li>Number of Transactions: " + strconv.Itoa(transactionCount) + "</li>"
	body += "<li>Number of Products Sold: " + strconv.Itoa(productSold) + "</li>"
	body += "<li>Total Income: " + strconv.Itoa(income) + "</li></ul>"
	body += "<p>Best regards,</p>"
	body += "<p>Kelontongpedia Team</p>"
	mail.SetBody("text/html", body)

	sender := gm.NewDialer("smtp.gmail.com", 25, "kelontongpedia2023@gmail.com", "xhdfxoyciurbfizb")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent to: ", email)
	}
}

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	if err := t.Execute(&buff, data); err != nil {
		return "", err
	}

	return buff.String(), nil
}
