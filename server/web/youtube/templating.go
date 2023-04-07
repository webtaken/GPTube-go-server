package web

import (
	"fmt"
	"log"
	"os"
	envManager "server/env_manager"
	"text/template"
)

type EmailTemplate struct {
	VideoID    string
	TotalCount int
	Votes1     int
	Votes2     int
	Votes3     int
	Votes4     int
	Votes5     int
}

func SendYoutubeTemplate(data EmailTemplate, subject string, emails []string) error {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tmplDir := fmt.Sprintf("%s%s", dir, "/web/templates/youtube/email_vanillacss.gotmpl")
	tmpl := template.Must(template.ParseFiles(tmplDir))
	newEmail := NewRequest(emails, subject, "")
	sendedData := struct {
		FrontendURL string
		Results     EmailTemplate
	}{
		FrontendURL: envManager.GoDotEnvVariable("FRONTEND_URL"),
		Results:     data,
	}
	err = newEmail.ParseTemplate(tmpl, sendedData)
	if err == nil {
		ok, err := newEmail.SendEmail()
		fmt.Println("Email sent: ", ok, err)
	} else {
		fmt.Println(err.Error())
	}

	return nil
}

func SendYoutubeErrorTemplate(subject string, emails []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	tmplDir := fmt.Sprintf("%s%s", dir, "/web/templates/youtube/error_vanillacss.gotmpl")
	tmpl := template.Must(template.ParseFiles(tmplDir))
	newEmail := NewRequest(emails, subject, "")

	err = newEmail.ParseTemplate(tmpl, nil)
	if err == nil {
		ok, err := newEmail.SendEmail()
		fmt.Println("Email error sent: ", ok, err)
	} else {
		fmt.Println(err.Error())
	}

	return nil
}
