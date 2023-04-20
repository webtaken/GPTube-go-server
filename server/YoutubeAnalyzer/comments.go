package YoutubeAnalyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	envManager "server/env_manager"
	"server/models"
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
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		return fmt.Errorf("unable to connect to the AI server")
	}

	return nil
}

func bertAnalysis(comments []*youtube.CommentThread, results *models.BertAIResults) error {
	// AI server information
	AIBertEndpoint := fmt.Sprintf(
		"%s/models/nlptown/bert-base-multilingual-uncased-sentiment",
		envManager.GoDotEnvVariable("AI_SERVER_URL"),
	)
	maxCharsAllowed := 512
	tmpResult := models.BertAIResults{}

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
			// Here we get the lowest comments //
			// pickup randomly the comments, we will take ~20% of the total comments
			// if math.random() <= 0.2 {
			// 	heap.push(comment)
			// }
			/////////////////////////////////////
			tmpResult.Score1++
		case "2 stars":
			tmpResult.Score2++
		case "3 stars":
			tmpResult.Score3++
		case "4 stars":
			tmpResult.Score4++
		default:
			tmpResult.Score5++
		}
	}

	for _, commentResults := range responseCommentsAI {
		getMaxScore(commentResults)
	}

	// Writing response to the global result
	mu.Lock()
	results.SuccessCount += len(responseCommentsAI)
	results.Score1 += tmpResult.Score1
	results.Score2 += tmpResult.Score2
	results.Score3 += tmpResult.Score3
	results.Score4 += tmpResult.Score4
	results.Score5 += tmpResult.Score5
	results.ErrorsCount += tmpResult.ErrorsCount
	mu.Unlock()

	return nil
}

func Analyze(body models.YoutubeAnalyzerReqBody) (*models.BertAIResults, error) {
	var part = []string{"id", "snippet"}
	results := &models.BertAIResults{}
	nextPageToken := ""

	// Check if AI server is running before calling Youtube API
	err := checkBertAIHealthcare()
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
			err = bertAnalysis(tmpComments, results)
			if err != nil {
				log.Printf("%v\n", err)
			}
		}()

		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	wg.Wait()
	return results, nil
}
