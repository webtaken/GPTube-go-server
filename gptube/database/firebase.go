package database

import (
	"context"
	"fmt"
	"gptube/config"
	"gptube/models"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var ctx context.Context
var sa option.ClientOption

func init() {
	ctx = context.Background()
	fmt.Printf("%s\n", ("ENV_MODE"))
	if config.Config("ENV_MODE") == "development" {
		sa = option.WithCredentialsFile("gptube-firebase-sdk-dev.json")
	} else {
		sa = option.WithCredentialsFile("gptube-firebase-sdk-prod.json")
	}
}

func AddYoutubeResult(results *models.YoutubeAnalyzerRespBody) (*firestore.DocumentRef, error) {
	collectionName := "YoutubeResults"
	// Use a service account
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	docRef, _, err := client.Collection(collectionName).Add(ctx, results)
	if err != nil {
		return nil, err
	}
	return docRef, nil
}
