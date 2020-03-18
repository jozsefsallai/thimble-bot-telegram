package utils

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Firebase will initialize a Firebase app
func Firebase() *firebase.App {
	options := option.WithCredentialsFile("./.firebase/service-account.json")
	app, err := firebase.NewApp(context.Background(), nil, options)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	return app
}
