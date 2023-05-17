package services

import (
	"container/heap"
	"fmt"
	"gptube/config"
	"gptube/models"
	"gptube/utils"
	"log"
	"math"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/youtube/v3"
)

var huggingFaceAuthHeader = fmt.Sprintf("Bearer %s", config.Config("HUGGING_FACE_TOKEN"))
var mu sync.Mutex
var AIEndpoints = map[string]string{
	"BERT": fmt.Sprintf(
		"%s/models/nlptown/bert-base-multilingual-uncased-sentiment",
		config.Config("AI_SERVER_URL"),
	),
	"RoBERTa": fmt.Sprintf(
		"%s/models/cardiffnlp/twitter-xlm-roberta-base-sentiment",
		config.Config("AI_SERVER_URL"),
	),
}

func CheckAIModelsWork() error {
	payload := []byte(`{"inputs":"i love you"}`)
	for _, endpoint := range AIEndpoints {
		agent := fiber.AcquireAgent()
		req := agent.Request()
		req.Header.Set("Authorization", huggingFaceAuthHeader)
		req.Header.Set("Content-Type", "application/json")
		req.Header.SetMethod(fiber.MethodPost)
		req.SetRequestURI(endpoint)
		agent.Body(payload)
		if err := agent.Parse(); err != nil {
			log.Println("[CheckAIModelsWork] error making the request: ", err)
			return err
		}
		code, _, errs := agent.String()
		if code != http.StatusOK && len(errs) > 0 {
			log.Println("[CheckAIModelsWork] error in response: ", errs[0])
			return errs[0]
		}
	}
	return nil
}

func MakeAICall(endpoint string, reqBody interface{}, resBody interface{}) error {
	agent := fiber.AcquireAgent()
	req := agent.Request()
	req.Header.Set("Authorization", huggingFaceAuthHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(endpoint)
	agent.JSON(reqBody)
	if err := agent.Parse(); err != nil {
		log.Println("[MakeAICall] error making the request: ", err)
		return err
	}
	code, _, errs := agent.Struct(resBody)
	if code != http.StatusOK && len(errs) > 0 {
		log.Println("[MakeAICall] error in response: ", errs[0])
		return errs[0]
	}
	return nil
}

func RobertaAnalysis(comments []*youtube.CommentThread, results *models.YoutubeAnalysisResults) error {
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

	err := MakeAICall(AIEndpoints["RoBERTa"], &reqRoberta, &resRoberta)
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

	err := MakeAICall(AIEndpoints["BERT"], &reqBert, &resBert)
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
