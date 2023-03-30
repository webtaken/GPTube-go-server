package YoutubeAnalyzer

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var Service *youtube.Service
var EnvVars map[string]string

func init() {
	env, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	EnvVars = env
	ctx := context.Background()
	Service, err = youtube.NewService(ctx, option.WithAPIKey(env["YOUTUBE_API_KEY"]))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
