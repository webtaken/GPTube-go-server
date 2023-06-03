package database

import (
	"container/heap"
	"context"
	"fmt"
	"gptube/config"
	"gptube/models"
	"log"

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

func AddYoutubeResult(results *models.YoutubeAnalyzerRespBody) error {
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatal("Couldn't start the firebase app.")
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}

	defer client.Close()
	resultsRef := client.Collection("YoutubeResults").Doc(results.VideoID)
	_, err = resultsRef.Set(ctx, results)
	if err != nil {
		return err
	}

	negativeCommentsRef := resultsRef.Collection("NegativeComments")
	for results.Results.NegativeComments.Len() > 0 {
		comment := heap.Pop(results.Results.NegativeComments).(*models.NegativeComment)
		_, _, err = negativeCommentsRef.Add(ctx, comment)
		if err != nil {
			log.Printf("Failed to add negative comment: %v", err)
		}
	}
	return nil
}
