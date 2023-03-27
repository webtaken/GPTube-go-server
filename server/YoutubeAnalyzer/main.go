package YoutubeAnalyzer

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var Service *youtube.Service

func init() {
	ctx := context.Background()
	var err error
	Service, err = youtube.NewService(ctx, option.WithAPIKey("AIzaSyAO6wG0jo7hMZ-FyLFIJli5X4TjfgQfM9o"))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
