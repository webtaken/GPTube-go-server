package YoutubeAnalyzer

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var Service *youtube.Service
var huggingFaceAuthHeader = fmt.Sprintf("Bearer %s",
	os.Getenv("HUGGING_FACE_TOKEN"))
var mu sync.Mutex

func init() {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	Service = service
}
