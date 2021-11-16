package main

import (
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

var (
	username = "******@163.com"
	host     = "smtp.163.com"
	passwd   = "******"
	to       = "******"
)

func main() {
	e := &email.Email{
		From:    "dfdf <" + username + ">",
		To:      []string{to},
		Subject: "3点几啦，点饭先～",
	}

	err := e.Send(host+":25", smtp.PlainAuth("", username, passwd, host))
	if err != nil {
		log.Println(err)
	}
}
