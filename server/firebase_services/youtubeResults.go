package firebase_services

import (
	"server/models"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

func AddYoutubeResult(results *models.YoutubeAnalyzerRespBody) (*firestore.DocumentRef, error) {
	collectionName := "youtubeResults"
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

	doc := make(map[string]interface{})
	doc["video_id"] = results.VideoID
	doc["video_title"] = results.VideoTitle
	doc["email"] = results.Email
	doc["error"] = results.Err
	doc["bert_results"] = map[string]interface{}{
		"score_1":       results.BertAnalysis.Score1,
		"score_2":       results.BertAnalysis.Score2,
		"score_3":       results.BertAnalysis.Score3,
		"score_4":       results.BertAnalysis.Score4,
		"score_5":       results.BertAnalysis.Score5,
		"errors_count":  results.BertAnalysis.ErrorsCount,
		"success_count": results.BertAnalysis.SuccessCount,
	}

	docRef, _, err := client.Collection(collectionName).Add(ctx, doc)
	if err != nil {
		return nil, err
	}
	return docRef, nil
}
