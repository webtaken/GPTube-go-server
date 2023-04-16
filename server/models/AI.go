package models

type BertAIResults struct {
	Score1       int `json:"score_1"`
	Score2       int `json:"score_2"`
	Score3       int `json:"score_3"`
	Score4       int `json:"score_4"`
	Score5       int `json:"score_5"`
	SuccessCount int `json:"success_count"`
	ErrorsCount  int `json:"errors_count"`
}

type YoutubeCommentsReqBertAI struct {
	Inputs []string `json:"inputs"`
}

type YoutubeCommentsResBertAI [][]struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}
