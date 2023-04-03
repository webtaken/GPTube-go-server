package web

import (
	"fmt"
	"log"
	"os"
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
	tmpl := template.Must(
		template.New("email_vanillacss.gotmpl").Funcs(
			template.FuncMap{
				"percentage": func(votes, totalCount int) int {
					if totalCount == 0 {
						return 0
					}
					return int(0.5 + (100 * float32(votes) / float32(totalCount)))
				},
			},
		).ParseFiles(tmplDir))
	newEmail := NewRequest([]string{"saul.rojas@ucsp.edu.pe"}, "GPTube Analysis", "Analysis ready")
	err = newEmail.ParseTemplate(tmpl, data)
	if err == nil {
		ok, err := newEmail.SendEmail()
		fmt.Println("Email sent: ", ok, err)
	} else {
		fmt.Println(err.Error())
	}

	return nil
}
