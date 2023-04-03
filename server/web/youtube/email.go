package web

import (
	"bytes"
	"net/smtp"
	envManager "server/env_manager"
	"text/template"
)

// Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		from:    envManager.GoDotEnvVariable("EMAIL_USER"),
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(
		addr, emailAuth, envManager.GoDotEnvVariable("EMAIL_USER"), r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(t *template.Template, data interface{}) error {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
