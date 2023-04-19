package firebase_services

import (
	"context"

	"google.golang.org/api/option"
)

var ctx context.Context
var sa option.ClientOption

func init() {
	ctx = context.Background()
	sa = option.WithCredentialsFile("gptube-firebase-sdk.json")
}
