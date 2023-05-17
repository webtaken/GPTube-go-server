package services

import (
	"bytes"
	"container/heap"
	"context"
	"fmt"
	"gptube/config"
	"gptube/models"
	"log"
	"os"
	"strconv"
	"sync"
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

	templateDirectory := fmt.Sprintf("%s%s", dir, "/templates/youtube_email_success.gotmpl")
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

	templateDirectory := fmt.Sprintf("%s%s", dir, "/templates/youtube_email_error.gotmpl")
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

func Analyze(body models.YoutubeAnalyzerReqBody) (*models.YoutubeAnalysisResults, error) {
	negativeComments := models.HeapNegativeComments([]*models.NegativeComment{})
	heap.Init(&negativeComments)
	limitComments := 10
	results := &models.YoutubeAnalysisResults{
		BertResults:           &models.BertAIResults{},
		RobertaResults:        &models.RobertaAIResults{},
		NegativeComments:      &negativeComments,
		NegativeCommentsLimit: limitComments,
	}

	var part = []string{"id", "snippet"}
	nextPageToken := ""

	// Check if AI services are running before calling Youtube API
	err := CheckAIModelsWork()

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	// Youtube calling
	call := Service.CommentThreads.List(part)
	call.VideoId(body.VideoID)
	for {
		if nextPageToken != "" {
			call.PageToken(nextPageToken)
		}
		response, err := call.Do()
		if err != nil {
			return results, err
		}

		tmpComments := make([]*youtube.CommentThread, len(response.Items))
		for i, p := range response.Items {
			if p == nil {
				continue
			}
			v := *p
			tmpComments[i] = &v
		}

		// Launching Two AI models to work in parallel
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = BertAnalysis(tmpComments, results)
			if err != nil {
				log.Printf("bert_analysis_error %v\n", err)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = RobertaAnalysis(tmpComments, results)
			if err != nil {
				log.Printf("bert_analysis_error %v\n", err)
			}
		}()
		//////////////////////////////////////////////

		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	wg.Wait()

	// Averaging results for roBERTa model
	results.RobertaResults.AverageResults()
	tmpHeap := models.HeapNegativeComments(make([]*models.NegativeComment, 0))
	for results.NegativeComments.Len() > 0 {
		item := heap.Pop(results.NegativeComments).(*models.NegativeComment)
		if tmpHeap.Len() <= results.NegativeCommentsLimit {
			heap.Push(&tmpHeap, item)
		} else {
			break
		}
	}
	results.NegativeComments = &tmpHeap

	return results, nil
}
