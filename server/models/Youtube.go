package models

import "google.golang.org/api/youtube/v3"

type YoutubePreAnalyzerReqBody struct {
	VideoID string `json:"video_id"`
}

type YoutubePreAnalyzerRespBody struct {
	VideoID       string                `json:"video_id,omitempty"`
	NumOfComments int                   `json:"number_of_comments,omitempty"`
	Snippet       *youtube.VideoSnippet `json:"snippet,omitempty"`
	Err           string                `json:"error,omitempty"`
}

type YoutubeAnalyzerReqBody struct {
	VideoID string `json:"video_id"`
	Email   string `json:"email,omitempty"`
	Err     string `json:"error,omitempty"`
}

type YoutubeCommentsReqBertAI struct {
	Inputs []string `json:"inputs"`
}

type YoutubeCommentsResBertAI [][]struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}
