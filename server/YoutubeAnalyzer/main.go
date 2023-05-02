package YoutubeAnalyzer

import (
	"context"
	"fmt"
	"log"
	envManager "server/env_manager"
	"sync"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var Service *youtube.Service
var huggingFaceAuthHeader = fmt.Sprintf("Bearer %s",
	envManager.GoDotEnvVariable("HUGGING_FACE_TOKEN"))
var mu sync.Mutex

const maxBadCommentsAllowed = 10

func init() {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(envManager.GoDotEnvVariable("YOUTUBE_API_KEY")))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	Service = service
}
