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

func SendTemplate(data EmailTemplate) error {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tmplDir := fmt.Sprintf("%s%s", dir, "/web/youtube/email_vanillacss.gotmpl")
	tmpl := template.Must(template.ParseFiles(tmplDir))
	newEmail := NewRequest([]string{"saul.rojas@ucsp.edu.pe"}, "GPTube Analysis", "Analysis ready")
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
