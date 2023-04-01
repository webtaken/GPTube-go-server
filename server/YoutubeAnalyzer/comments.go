package YoutubeAnalyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	envManager "server/env_manager"
	"server/models"
	"sync"

	strip "github.com/grokify/html-strip-tags-go"
	"google.golang.org/api/youtube/v3"
)

func bertAnalysis(comments []*youtube.CommentThread, w *sync.WaitGroup) {
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
		log.Fatal(err)
	}

	response, err := http.Post(AIBertEndpoint, "application/json",
		bytes.NewBuffer(jsonCommentsForAI))

	if err != nil {
		log.Fatal(err)
	}

	responseCommentsForAI := make([]models.YoutubeCommentThreadForAI, 0)
	json.NewDecoder(response.Body).Decode(&responseCommentsForAI)
	for _, comment := range responseCommentsForAI {
		fmt.Printf("id: %s\tsentiment score:%d\n", comment.CommentID, comment.SentimentScore)
	}
	w.Done()
}

func GetComments(youtubeRequestBody models.YoutubeAnalyzerRequestBody) ([]*youtube.CommentThread, error) {
	var part = []string{"id", "snippet"}
	comments := make([]*youtube.CommentThread, 0)
	nextPageToken := ""

	var wg sync.WaitGroup
	for {
		call := Service.CommentThreads.List(part)
		call.VideoId(youtubeRequestBody.VideoID)
		if nextPageToken != "" {
			call.PageToken(nextPageToken)
		}
		response, err := call.Do()
		if err != nil {
			return comments, err
		}

		// Here goes AI analysis with goroutines //
		wg.Add(1)
		go bertAnalysis(response.Items, &wg)
		///////////////////////////////////////////

		comments = append(comments, response.Items...)

		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	wg.Wait()
	return comments, nil
}
