package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/florian-lahitte-uvi/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
		return
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("emailTemplate/%s.html", m.Template))
		if err != nil {
			errorLog.Println(err)
			return
		}
		mailTemplate := string(data)
		msgTosend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgTosend)
	}
	err = email.Send(client)
	if err != nil {
		errorLog.Println(err)
		return
	} else {
		infoLog.Println("Mail sent!")
	}
}
