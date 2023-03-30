package YoutubeAnalyzer

import (
	"server/models"

	"google.golang.org/api/youtube/v3"
)

func GetComments(youtubeRequestBody models.YoutubeAnalyzerRequestBody) ([]*youtube.CommentThread, error) {
	var part = []string{"id", "snippet"}
	comments := make([]*youtube.CommentThread, 0)
	nextPageToken := ""
	for {
		call := Service.CommentThreads.List(part)
		call.VideoId(youtubeRequestBody.VideoID)
		if nextPageToken != "" {
			call.PageToken(nextPageToken)
		}
		response, err := call.Do()
		if err != nil {
			return comments, err
		}

		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
		comments = append(comments, response.Items...)
	}
	return comments, nil
}
