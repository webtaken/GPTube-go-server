package firebase_services

import (
	"context"
	"fmt"
	"gptube/environment"

	"google.golang.org/api/option"
)

var ctx context.Context
var sa option.ClientOption

func init() {
	ctx = context.Background()
	fmt.Printf("%s\n", ("ENV_MODE"))
	if environment.Getenv("ENV_MODE") == "development" {
		sa = option.WithCredentialsFile("gptube-firebase-sdk-dev.json")
	} else {
		sa = option.WithCredentialsFile("gptube-firebase-sdk-prod.json")
	}
}
