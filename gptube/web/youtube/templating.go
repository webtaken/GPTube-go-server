package web

import (
	"fmt"
	"gptube/models"
	"log"
	"os"
	"text/template"
)

func SendYoutubeTemplate(data models.YoutubeAnalyzerRespBody, subject string, emails []string) error {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	templateDirectory := fmt.Sprintf("%s%s", dir, "/web/templates/youtube/email_success.gotmpl")
	template := template.Must(template.ParseFiles(templateDirectory))
	newEmail := NewRequest(emails, subject, "")
	sendedData := struct {
		FrontendURL string
		Results     models.YoutubeAnalyzerRespBody
	}{
		FrontendURL: os.Getenv("FRONTEND_URL"),
		Results:     data,
	}
	err = newEmail.ParseTemplate(template, sendedData)
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

	templateDirectory := fmt.Sprintf("%s%s", dir, "/web/templates/youtube/email_error.gotmpl")
	template := template.Must(template.ParseFiles(templateDirectory))
	newEmail := NewRequest(emails, subject, "")

	err = newEmail.ParseTemplate(template, nil)
	if err == nil {
		ok, err := newEmail.SendEmail()
		fmt.Println("Email error sent: ", ok, err)
	} else {
		fmt.Println(err.Error())
	}

	return nil
}
