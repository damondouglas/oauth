package firestore

import (
	"log"

	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"
)

func firebaseApp(pathToCredentials string) (*firebase.App, error) {

	opt := option.WithCredentialsFile(pathToCredentials)
	return firebase.NewApp(context.Background(), nil, opt)
}
