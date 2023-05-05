package firebase_services

import (
	"context"
	envManager "server/env_manager"

	"google.golang.org/api/option"
)

var ctx context.Context
var sa option.ClientOption

func init() {
	ctx = context.Background()
	if envManager.GoDotEnvVariable("ENV_MODE") == "development" {
		sa = option.WithCredentialsFile("gptube-firebase-sdk-dev.json")
	} else {
		sa = option.WithCredentialsFile("gptube-firebase-sdk-prod.json")
	}
}
