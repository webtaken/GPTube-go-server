package web

import (
	"net/smtp"
	envManager "server/env_manager"
)

var emailAuth smtp.Auth

func init() {
	emailAuth = smtp.PlainAuth(
		"",
		envManager.GoDotEnvVariable("EMAIL_USER"),
		envManager.GoDotEnvVariable("EMAIL_PASSWORD"),
		"smtp.gmail.com",
	)
}
