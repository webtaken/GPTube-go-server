package YoutubeAnalyzer

import (
	"fmt"
	envManager "server/env_manager"
	"server/models"
	"strconv"

	"google.golang.org/api/youtube/v3"
)

func CanProcessVideo(youtubeRequestBody *models.YoutubePreAnalyzerReqBody) (*youtube.VideoListResponse, error) {
	// The max number of comments we can process
	maxNumberOfComments, _ := strconv.Atoi(envManager.GoDotEnvVariable("YOUTUBE_MAX_COMMENTS_CAPACITY"))

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
