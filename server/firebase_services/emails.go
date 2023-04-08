package firebase_services

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func AddLead(email string) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("gptube-firebase-sdk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	defer client.Close()

	_, _, err = client.Collection("leads").Add(ctx, map[string]interface{}{
		"email": email,
	})
	if err != nil {
		log.Printf("Failed adding email %s: %v", email, err)
	}
}
