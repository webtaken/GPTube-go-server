package firebase_services

import (
	"gptube/models"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

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
