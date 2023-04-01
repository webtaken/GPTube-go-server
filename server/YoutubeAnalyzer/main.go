package YoutubeAnalyzer

import (
	"context"
	"log"
	envManager "server/env_manager"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var Service *youtube.Service

func init() {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(envManager.GoDotEnvVariable("YOUTUBE_API_KEY")))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	Service = service
}
