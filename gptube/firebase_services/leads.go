package firebase_services

import (
	firebase "firebase.google.com/go"
)

func AddLead(email string) error {
	collectionName := "leads"
	// Use a service account
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}

	defer client.Close()

	_, _, err = client.Collection(collectionName).Add(ctx, map[string]interface{}{
		"email": email,
	})
	if err != nil {
		return err
	}
	return nil
}
