package YoutubeAnalyzer

import (
	"fmt"
	"server/models"

	"google.golang.org/api/youtube/v3"
)

func CanProcessVideo(youtubeRequestBody models.YoutubeAnalyzerRequestBody) (*youtube.VideoListResponse, error) {
	// The max number of comments we can process
	maxNumberOfComments := 8000

	var part = []string{"snippet", "contentDetails", "statistics"}
	call := Service.Videos.List(part)
	call.Id(youtubeRequestBody.VideoID)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	if len(response.Items) == 0 {
		return nil, fmt.Errorf("video not found")
	} else if response.Items[0].Statistics.CommentCount > uint64(maxNumberOfComments) {
		return nil, fmt.Errorf("max number of comments to process exceeded")
	}
	return response, nil
}
