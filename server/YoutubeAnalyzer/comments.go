package YoutubeAnalyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	envManager "server/env_manager"
	"server/models"
	web "server/web/youtube"
	"sync"

	strip "github.com/grokify/html-strip-tags-go"
	"google.golang.org/api/youtube/v3"
)

func checkBertAIHealthcare() error {
	// AI server information
	AIBertEndpoint := fmt.Sprintf(
		"%s/models/nlptown/bert-base-multilingual-uncased-sentiment",
		envManager.GoDotEnvVariable("AI_SERVER_URL"),
	)

	payload := []byte(`{"inputs":"i love you"}`)
	client := &http.Client{}
	req, err := http.NewRequest("POST", AIBertEndpoint, bytes.NewBuffer(payload))

	if err != nil {
		log.Println("Error creating request: ", err)
		return err
	}

	req.Header.Set("Authorization", huggingFaceAuthHeader)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making request: ", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to connect to the AI server")
	}

	return nil
}

func bertAnalysis(comments []*youtube.CommentThread, results *web.EmailTemplate) error {
	// AI server information
	AIBertEndpoint := fmt.Sprintf(
		"%s/models/nlptown/bert-base-multilingual-uncased-sentiment",
		envManager.GoDotEnvVariable("AI_SERVER_URL"),
	)
	maxCharsAllowed := 512
	tmpResult := web.EmailTemplate{}

	requestCommentsAI := models.YoutubeCommentsReqBertAI{Inputs: make([]string, 0)}
	for _, comment := range comments {
		cleanComment := strip.StripTags(comment.Snippet.TopLevelComment.Snippet.TextDisplay)
		if len(cleanComment) <= maxCharsAllowed {
			requestCommentsAI.Inputs = append(requestCommentsAI.Inputs, cleanComment)
		} else {
			tmpResult.ErrorsCount++
		}
	}

	// Here goes the Call to BERT model in the AI API
	jsonRequestCommentsAI, err := json.Marshal(requestCommentsAI)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", AIBertEndpoint, bytes.NewBuffer(jsonRequestCommentsAI))

	if err != nil {
		log.Println("Error creating request: ", err)
		return err
	}

	req.Header.Set("Authorization", huggingFaceAuthHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making request: ", err)
		return err
	}

	defer resp.Body.Close()

	responseCommentsAI := models.YoutubeCommentsResBertAI{}
	err = json.NewDecoder(resp.Body).Decode(&responseCommentsAI)
	if err != nil {
		return err
	}

	getMaxScore := func(commentResults []struct {
		Label string  "json:\"label\""
		Score float64 "json:\"score\""
	}) {
		tmpBertScore := struct {
			Label string  "json:\"label\""
			Score float64 "json:\"score\""
		}{
			Label: "5 stars",
			Score: math.Inf(-1),
		}

		for _, bertScore := range commentResults {
			if bertScore.Score > tmpBertScore.Score {
				tmpBertScore.Label = bertScore.Label
				tmpBertScore.Score = bertScore.Score
			}
		}

		switch tmpBertScore.Label {
		case "1 star":
			tmpResult.Votes1++
		case "2 stars":
			tmpResult.Votes2++
		case "3 stars":
			tmpResult.Votes3++
		case "4 stars":
			tmpResult.Votes4++
		default:
			tmpResult.Votes5++
		}
	}

	for _, commentResults := range responseCommentsAI {
		getMaxScore(commentResults)
	}

	// Writing response to the global result
	mu.Lock()
	results.TotalCount += len(responseCommentsAI)
	results.Votes1 += tmpResult.Votes1
	results.Votes2 += tmpResult.Votes2
	results.Votes3 += tmpResult.Votes3
	results.Votes4 += tmpResult.Votes4
	results.Votes5 += tmpResult.Votes5
	results.ErrorsCount += tmpResult.ErrorsCount
	mu.Unlock()

	return nil
}

func GetComments(youtubeRequestBody models.YoutubeAnalyzerRequestBody) (*web.EmailTemplate, error) {
	var part = []string{"id", "snippet"}
	commentsResults := &web.EmailTemplate{}
	nextPageToken := ""

	// Check if AI server is running before calling Youtube API
	err := checkBertAIHealthcare()
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

		// Copy extracted comments to another variable to send to the analysis tool
		// This implementation avoids race conditions while analyzing comments extracted
		// on "response.Items" variables
		tmpComments := make([]*youtube.CommentThread, len(response.Items))
		// Copy every Item to tmpComments to avoid race condition
		for i, p := range response.Items {
			if p == nil {
				// Skip to next for nil source pointer
				continue
			}
			// Create shallow copy of source element
			v := *p
			// Assign address of copy to destination.
			tmpComments[i] = &v
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err = bertAnalysis(tmpComments, commentsResults)
			if err != nil {
				log.Printf("%v\n", err)
			}
		}()

		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	// for _, page := range pages {
	// 	wg.Add(1)
	// 	go func(pageToken string) {
	// 		defer wg.Done()
	// 		newCall := Service.CommentThreads.List(part)
	// 		newCall.VideoId(youtubeRequestBody.VideoID)
	// 		newCall.PageToken(pageToken)
	// 		response, err := newCall.Do()
	// 		if err != nil {
	// 			return
	// 		}
	// 		err = bertAnalysis(response.Items, commentsResults)
	// 		if err != nil {
	// 			log.Printf("%v\n", err)
	// 		}
	// 	}(page)
	// }
	wg.Wait()
	return commentsResults, nil
}
