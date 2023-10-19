package main

import (
	"fmt"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

type SMTPServer struct {
	// host     string
	// port     int
	// username string
	// password string

	d *gomail.Dialer
}

type SMTPRequestBody struct {
	Subject string   `json:"subject"`
	To      []string `json:"to"`
	From    string   `json:"from"`
	Body    string   `json:"body"`
}

func NewSmtpServer() *SMTPServer {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	d := gomail.NewDialer(host, port, username, password)
	d.StartTLSPolicy = gomail.MandatoryStartTLS

	return &SMTPServer{
		d: d,
	}
}

func (ss SMTPServer) SendMail(smtpBody SMTPRequestBody) error {
	m := gomail.NewMessage()
	fmt.Printf("smtpBody: %v\n", smtpBody)
	m.SetHeader("From", smtpBody.From)
	m.SetHeader("To", smtpBody.To...)
	m.SetHeader("Subject", smtpBody.Subject)

	m.SetBody("text/plain", smtpBody.Body)

	if err := ss.d.DialAndSend(m); err != nil {
		panic(err)
	}
	return nil
}
