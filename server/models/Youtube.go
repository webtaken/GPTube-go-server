package models

type YoutubeAnalyzerRequestBody struct {
	VideoID string `json:"video_id"`
}

type YoutubeCommentThreadForAI struct {
	CommentID      string `json:"comment_id"`
	TextDisplay    string `json:"text_display"`
	SentimentScore int    `json:"sentiment_score,omitempty"`
}
