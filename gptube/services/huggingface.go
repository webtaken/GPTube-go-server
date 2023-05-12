package services

import (
	"bytes"
	"container/heap"
	"encoding/json"
	"fmt"
	"gptube/config"
	"gptube/models"
	"gptube/utils"
	"io"
	"log"
	"math"
	"net/http"
	"sync"

	"google.golang.org/api/youtube/v3"
)

var huggingFaceAuthHeader = fmt.Sprintf("Bearer %s", config.Config("HUGGING_FACE_TOKEN"))
var mu sync.Mutex

func CheckAIModelsWork() error {
	AIBertEndpoint := fmt.Sprintf(
		"%s/models/nlptown/bert-base-multilingual-uncased-sentiment",
		config.Config("AI_SERVER_URL"),
	)
	AIRobertaEndpoint := fmt.Sprintf(
		"%s/models/cardiffnlp/twitter-xlm-roberta-base-sentiment",
		config.Config("AI_SERVER_URL"),
	)
	var AIEndpoints = []string{AIBertEndpoint, AIRobertaEndpoint}

	payload := []byte(`{"inputs":"i love you"}`)
	client := &http.Client{}

	for _, endpoint := range AIEndpoints {
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))

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
			return fmt.Errorf("unable to connect to the AI services")
		}
	}
	return nil
}

func MakeAICall(endpoint string, reqBody interface{}, resBody interface{}) error {
	// Here goes the Call to BERT model in the AI API
	jsonReqAI, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonReqAI))

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

	if resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("some strings have incorrect format")
	}

	err = json.NewDecoder(resp.Body).Decode(resBody)
	if err != nil {
		log.Printf("%v", resp.Body)
		return err
	}

	return nil
}

func RobertaAnalysis(comments []*youtube.CommentThread, results *models.YoutubeAnalysisResults) error {
	RobertaEndpoint := fmt.Sprintf(
		"%s/models/cardiffnlp/twitter-xlm-roberta-base-sentiment",
		config.Config("AI_SERVER_URL"),
	)

	tmpResults := &models.RobertaAIResults{}
	reqRoberta := models.ReqRobertaAI{Inputs: make([]string, 0)}
	resRoberta := models.ResRobertaAI{}

	validComments := make([]*youtube.Comment, 0)
	maxCharsAllow := 512
	for _, comment := range comments {
		clean := utils.CleanComment(comment.Snippet.TopLevelComment.Snippet.TextOriginal)
		if len(clean) <= maxCharsAllow {
			reqRoberta.Inputs = append(reqRoberta.Inputs, clean)
			validComments = append(validComments, comment.Snippet.TopLevelComment)
		} else {
			tmpResults.ErrorsCount++
		}
	}

	err := MakeAICall(RobertaEndpoint, &reqRoberta, &resRoberta)
	if err != nil {
		return err
	}

	tmpResults.SuccessCount = len(resRoberta)

	negativeComments := make(models.HeapNegativeComments, 0)
	heap.Init(&negativeComments)

	callback := func(commentResults models.ResAISchema, comment *youtube.Comment) {
		for _, result := range commentResults {
			switch result.Label {
			case "negative":
				badComment := models.Comment{
					CommentID:             comment.Id,
					TextDisplay:           comment.Snippet.TextDisplay,
					TextOriginal:          comment.Snippet.TextOriginal,
					AuthorDisplayName:     comment.Snippet.AuthorDisplayName,
					AuthorProfileImageUrl: comment.Snippet.AuthorProfileImageUrl,
					ParentID:              comment.Snippet.ParentId,
					LikeCount:             comment.Snippet.LikeCount,
					ModerationStatus:      comment.Snippet.ModerationStatus,
				}
				item := &models.NegativeComment{
					Comment:  &badComment,
					Priority: result.Score,
				}
				heap.Push(&negativeComments, item)
				tmpResults.Negative += result.Score
			case "neutral":
				tmpResults.Neutral += result.Score
			default:
				tmpResults.Positive += result.Score
			}
		}
	}

	for i := 0; i < len(resRoberta); i++ {
		callback(resRoberta[i], validComments[i])
	}

	// Writing response to the global result
	mu.Lock()
	results.RobertaResults.Positive += tmpResults.Positive
	results.RobertaResults.Negative += tmpResults.Negative
	results.RobertaResults.Neutral += tmpResults.Neutral
	results.RobertaResults.ErrorsCount += tmpResults.ErrorsCount
	results.RobertaResults.SuccessCount += tmpResults.SuccessCount

	// Most negative comments from roBERTa model
	for negativeComments.Len() > 0 {
		item := heap.Pop(&negativeComments).(*models.NegativeComment)
		heap.Push(results.NegativeComments, item)
	}
	mu.Unlock()

	return nil
}

func BertAnalysis(comments []*youtube.CommentThread, results *models.YoutubeAnalysisResults) error {
	BertEndpoint := fmt.Sprintf(
		"%s/models/nlptown/bert-base-multilingual-uncased-sentiment",
		config.Config("AI_SERVER_URL"),
	)

	tmpResults := &models.BertAIResults{}
	reqBert := models.ReqBertAI{Inputs: make([]string, 0)}
	resBert := models.ResBertAI{}

	maxCharsAllow := 512
	for _, comment := range comments {
		clean := utils.CleanComment(comment.Snippet.TopLevelComment.Snippet.TextOriginal)
		if len(clean) <= maxCharsAllow {
			reqBert.Inputs = append(reqBert.Inputs, clean)
		} else {
			tmpResults.ErrorsCount++
		}
	}

	err := MakeAICall(BertEndpoint, &reqBert, &resBert)
	if err != nil {
		return err
	}

	tmpResults.SuccessCount = len(resBert)

	callback := func(commentResults models.ResAISchema) {
		tmpBertScore := models.ResAISchema{{
			Label: "1 star",
			Score: math.Inf(-1),
		}}[0]

		for _, result := range commentResults {
			if result.Score > tmpBertScore.Score {
				tmpBertScore.Label = result.Label
				tmpBertScore.Score = result.Score
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
			tmpResults.Score1++
		case "2 stars":
			tmpResults.Score2++
		case "3 stars":
			tmpResults.Score3++
		case "4 stars":
			tmpResults.Score4++
		default:
			tmpResults.Score5++
		}
	}

	for i := 0; i < len(resBert); i++ {
		callback(resBert[i])
	}

	// Writing response to the global result
	mu.Lock()
	results.BertResults.Score1 += tmpResults.Score1
	results.BertResults.Score2 += tmpResults.Score2
	results.BertResults.Score3 += tmpResults.Score3
	results.BertResults.Score4 += tmpResults.Score4
	results.BertResults.Score5 += tmpResults.Score5
	results.BertResults.ErrorsCount += tmpResults.ErrorsCount
	results.BertResults.SuccessCount += tmpResults.SuccessCount
	mu.Unlock()

	return nil
}
