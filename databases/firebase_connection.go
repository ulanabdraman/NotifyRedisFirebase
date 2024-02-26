// firebase_connection.go
package databases

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func ConnectToFirebase() (*firebase.App, error) {
	opt := option.WithCredentialsFile("glot-8e244-firebase-adminsdk-3vt2y-3eada6f222.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal(err)
	}

	return app, err
}
