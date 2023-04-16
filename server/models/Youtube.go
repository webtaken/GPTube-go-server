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
	VideoID    string `json:"video_id,omitempty"`
	VideoTitle string `json:"video_title,omitempty"`
	Email      string `json:"email,omitempty"`
}

type YoutubeAnalyzerRespBody struct {
	VideoID      string         `json:"video_id"`
	Email        string         `json:"email,omitempty"`
	BertAnalysis *BertAIResults `json:"bert_analysis,omitempty"`
	ResultsID    string         `json:"results_id,omitempty"` // This is the _id for results in the fireStore database
	Err          string         `json:"error,omitempty"`
}
