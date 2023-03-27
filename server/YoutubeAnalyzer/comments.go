package YoutubeAnalyzer

import (
	"fmt"
	"log"
	"server/models"
)

func getComments(youtubeRequestBody models.YoutubeAnalyzerRequestBody) {
	var part = []string{"id", "snippet"}
	call := Service.CommentThreads.List(part)
	call.VideoId(youtubeRequestBody.VideoID)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	nextPage := response.NextPageToken
	fmt.Printf("NEXT PAGE TOKEN: %s\n\n", nextPage)
	for _, comment := range response.Items {
		commentID := comment.Id
		originalCommentText := comment.Snippet.TopLevelComment.Snippet.TextOriginal

		// Print the comment ID and original text for the commentThreads resource.
		fmt.Println(commentID, ": ", originalCommentText)
	}
}
