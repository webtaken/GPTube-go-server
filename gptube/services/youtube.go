package services

import (
	"bytes"
	"context"
	"fmt"
	"gptube/config"
	"gptube/models"
	"log"
	"os"
	"strconv"
	"text/template"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var Service *youtube.Service

func init() {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(config.Config("YOUTUBE_API_KEY")))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	Service = service
}

func (r *Request) ParseTemplate(t *template.Template, data interface{}) error {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

func SendYoutubeSuccessTemplate(data models.YoutubeAnalyzerRespBody, subject string, emails []string) error {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	templateDirectory := fmt.Sprintf("%s%s", dir, "/templates/email_success.gotmpl")
	template := template.Must(template.ParseFiles(templateDirectory))
	newEmail := NewRequest(emails, subject, "")
	sendedData := struct {
		FrontendURL string
		Results     models.YoutubeAnalyzerRespBody
	}{
		FrontendURL: config.Config("FRONTEND_URL"),
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

func CanProcessVideo(youtubeRequestBody *models.YoutubePreAnalyzerReqBody) (*youtube.VideoListResponse, error) {
	// The max number of comments we can process
	maxNumberOfComments, _ := strconv.Atoi(config.Config("YOUTUBE_MAX_COMMENTS_CAPACITY"))

	var part = []string{"snippet", "contentDetails", "statistics"}
	call := Service.Videos.List(part)
	call.Id(youtubeRequestBody.VideoID)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	if len(response.Items) == 0 {
		return nil, fmt.Errorf("video not found")
	} else if response.Items[0].Statistics.CommentCount > uint64(maxNumberOfComments) {
		return nil, fmt.Errorf("max number of comments to process exceeded")
	}
	return response, nil
}
