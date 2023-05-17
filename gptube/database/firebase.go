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
	fmt.Printf("In %s mode.\n", config.Config("ENV_MODE"))
	if config.Config("ENV_MODE") == "development" {
		fmt.Printf("Starting the firebase sdk: %s\n", config.Config("DB_KEYS_DEVELOPMENT"))
		sa = option.WithCredentialsFile(config.Config("DB_KEYS_DEVELOPMENT"))
	} else {
		fmt.Printf("Starting the firebase sdk: %s\n", config.Config("DB_KEYS_PRODUCTION"))
		sa = option.WithCredentialsFile(config.Config("DB_KEYS_PRODUCTION"))
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
