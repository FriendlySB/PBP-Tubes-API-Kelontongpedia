package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"bytes"
	"fmt"
	"text/template"

	gm "gopkg.in/gomail.v2"
)

func sendMail(user model.User) {
	mail := gm.NewMessage()

	template := "bin/template/mail.html"

	result, _ := parseTemplate(template, user)

	mail.SetHeader("From", "lamabunta@gmail.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Notifications")
	mail.SetBody("text/html", result)

	sender := gm.NewDialer("smtp.gmail.com", 25, "lamabunta@gmail.com", "gnkglansnfmbshty")
	
	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent to: ", user.Email)
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
