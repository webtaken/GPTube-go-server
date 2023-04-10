package models

type YoutubeAnalyzerRequestBody struct {
	VideoID string `json:"video_id"`
	Email   string `json:"email"`
}

type YoutubeCommentsReqBertAI struct {
	Inputs []string `json:"inputs"`
}

type YoutubeCommentsResBertAI [][]struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}
