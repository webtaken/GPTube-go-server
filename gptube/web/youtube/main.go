package web

import (
	"net/smtp"
	"os"
)

var emailAuth smtp.Auth

func init() {
	emailAuth = smtp.PlainAuth(
		"",
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASSWORD"),
		"smtp.gmail.com",
	)
}
