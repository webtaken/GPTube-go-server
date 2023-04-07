package YoutubeAnalyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	envManager "server/env_manager"
	"server/models"
	web "server/web/youtube"
	"sync"

	strip "github.com/grokify/html-strip-tags-go"
	"google.golang.org/api/youtube/v3"
)

func checkAIServerHealthcare() error {
	resp, err := http.Get(envManager.GoDotEnvVariable("AI_SERVER_URL"))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func bertAnalysis(comments []*youtube.CommentThread, results *web.EmailTemplate) error {
	// AI server information
	AIBertEndpoint := fmt.Sprintf("%s/YT", envManager.GoDotEnvVariable("AI_SERVER_URL"))
	commentsForAI := make([]models.YoutubeCommentThreadForAI, 0)
	for _, comment := range comments {
		commentsForAI = append(commentsForAI, models.YoutubeCommentThreadForAI{
			CommentID: comment.Id,
			// avoid html tags inside text
			TextDisplay:    strip.StripTags(comment.Snippet.TopLevelComment.Snippet.TextDisplay),
			SentimentScore: 0,
		})
	}

	// Here goes the Call to BERT model in the AI API
	jsonCommentsForAI, err := json.Marshal(commentsForAI)
	if err != nil {
		return err
	}

	response, err := http.Post(AIBertEndpoint, "application/json",
		bytes.NewBuffer(jsonCommentsForAI))

	if err != nil {
		return err
	}

	responseCommentsForAI := make([]models.YoutubeCommentThreadForAI, 0)
	err = json.NewDecoder(response.Body).Decode(&responseCommentsForAI)
	if err != nil {
		return err
	}

	tmpResult := web.EmailTemplate{}
	for _, comment := range responseCommentsForAI {
		switch comment.SentimentScore {
		case 1:
			tmpResult.Votes1++
		case 2:
			tmpResult.Votes2++
		case 3:
			tmpResult.Votes3++
		case 4:
			tmpResult.Votes4++
		default:
			tmpResult.Votes5++
		}
	}

	// Writing response to the global result
	mu.Lock()
	results.TotalCount += len(responseCommentsForAI)
	results.Votes1 += tmpResult.Votes1
	results.Votes2 += tmpResult.Votes2
	results.Votes3 += tmpResult.Votes3
	results.Votes4 += tmpResult.Votes4
	results.Votes5 += tmpResult.Votes5
	mu.Unlock()

	return nil
}

func GetComments(youtubeRequestBody models.YoutubeAnalyzerRequestBody) (*web.EmailTemplate, error) {

	var part = []string{"id", "snippet"}
	commentsResults := &web.EmailTemplate{}
	pages := []string{""}
	nextPageToken := ""

	// Check if AI server is running before calling Youtube API
	err := checkAIServerHealthcare()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	// Youtube calling
	call := Service.CommentThreads.List(part)
	call.VideoId(youtubeRequestBody.VideoID)

	for {
		if nextPageToken != "" {
			call.PageToken(nextPageToken)
		}
		response, err := call.Do()
		if err != nil {
			return commentsResults, err
		}
		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
		pages = append(pages, nextPageToken)
	}

	for _, page := range pages {
		wg.Add(1)
		go func(pageToken string) {
			defer wg.Done()
			newCall := Service.CommentThreads.List(part)
			newCall.VideoId(youtubeRequestBody.VideoID)
			newCall.PageToken(pageToken)
			response, err := newCall.Do()
			if err != nil {
				return
			}
			err = bertAnalysis(response.Items, commentsResults)
			if err != nil {
				log.Printf("%v\n", err)
			}
		}(page)
	}
	wg.Wait()
	return commentsResults, nil
}
