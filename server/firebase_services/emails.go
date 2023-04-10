package firebase_services

import (
	"context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func AddLead(email string) error {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("gptube-firebase-sdk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}

	defer client.Close()

	_, _, err = client.Collection("leads").Add(ctx, map[string]interface{}{
		"email": email,
	})
	if err != nil {
		return err
	}
	return nil
}
