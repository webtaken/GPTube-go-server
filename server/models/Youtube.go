package models

import "google.golang.org/api/youtube/v3"

type Comment struct {
	CommentID             string `json:"commentID,omitempty" firestore:"commentID,omitempty"`
	TextDisplay           string `json:"textDisplay,omitempty" firestore:"textDisplay,omitempty"`
	TextOriginal          string `json:"textOriginal,omitempty" firestore:"textOriginal,omitempty"`
	AuthorDisplayName     string `json:"authorDisplayName,omitempty" firestore:"authorDisplayName,omitempty"`
	AuthorProfileImageUrl string `json:"authorProfileImageUrl,omitempty" firestore:"authorProfileImageUrl,omitempty"`
	ParentID              string `json:"parentID,omitempty" firestore:"parentID,omitempty"`
	LikeCount             int64  `json:"likeCount,omitempty" firestore:"likeCount,omitempty"`
	// ModerationStatus: The comment's moderation status. Will not be set if
	// the comments were requested through the id filter.
	//
	// Possible values:
	//   "published" - The comment is available for public display.
	//   "heldForReview" - The comment is awaiting review by a moderator.
	//   "likelySpam"
	//   "rejected" - The comment is unfit for display.
	ModerationStatus string `json:"moderationStatus,omitempty" firestore:"moderationStatus,omitempty"`
}

type YoutubePreAnalyzerReqBody struct {
	VideoID string `json:"video_id"`
}

type YoutubePreAnalyzerRespBody struct {
	VideoID       string                `json:"video_id,omitempty"`
	NumOfComments int                   `json:"number_of_comments,omitempty"`
	RequiresEmail bool                  `json:"requires_email,omitempty"`
	Snippet       *youtube.VideoSnippet `json:"snippet,omitempty"`
	Err           string                `json:"error,omitempty"`
}

type YoutubeAnalyzerReqBody struct {
	VideoID    string `json:"video_id,omitempty"`
	VideoTitle string `json:"video_title,omitempty"`
	Email      string `json:"email,omitempty"`
}

type YoutubeAnalyzerRespBody struct {
	VideoID    string                  `json:"video_id" firestore:"video_id"`
	VideoTitle string                  `json:"video_title,omitempty" firestore:"video_title,omitempty"`
	Email      string                  `json:"email,omitempty" firestore:"email,omitempty"`
	Results    *YoutubeAnalysisResults `json:"-" firestore:"results,omitempty"`
	// This is the _id for results in the fireStore database
	ResultsID string `json:"results_id,omitempty" firestore:"-"`
	// Errors if encountered
	Err string `json:"error,omitempty" firestore:"error,omitempty"`
}

type YoutubeAnalysisResults struct {
	VideoID               string                `json:"video_id,omitempty" firestore:"video_id,omitempty"`
	VideoTitle            string                `json:"video_title,omitempty" firestore:"video_title,omitempty"`
	BertResults           *BertAIResults        `json:"bert_results,omitempty" firestore:"bert_results,omitempty"`
	RobertaResults        *RobertaAIResults     `json:"roberta_results,omitempty" firestore:"roberta_results,omitempty"`
	NegativeComments      *HeapNegativeComments `json:"negative_comments,omitempty" firestore:"negative_comments,omitempty"`
	NegativeCommentsLimit int                   `json:"-" firestore:"-"`
}
