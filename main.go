package main

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	projectID := "my-first-project-3c6af"
	serviceAccountKeyPath := "./key.json"

	// Initialize the Firebase app
	ctx := context.Background()
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: projectID,
	}, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v\n", err)
	}

	// Initialize the Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v\n", err)
	}
	defer client.Close()

	collectionName := "events"

	startDate := time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 3, 8, 0, 0, 0, 0, time.UTC)

	var solvedValues []int64

	// Iterate over the date range and retrieve the documents with that date as the name
	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		docRef := client.Collection(collectionName).Doc(date.Format("2006-01-02"))
		doc, err := docRef.Get(ctx)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				fmt.Printf("Document with ID %s does not exist\n", docRef.ID)
			} else {
				log.Fatalf("Failed to retrieve document %s: %v\n", docRef.ID, err)
			}
		} else {
			data := doc.Data()
			if val, ok := data["solved"].(int64); ok {
				// fmt.Printf("Document %s has field 'solved' with value %d\n", docRef.ID, val)
				solvedValues = append(solvedValues, val)
			} else {
				fmt.Printf("Document %s does not have field 'solved'\n", docRef.ID)
			}
		}
	}
	fmt.Println(solvedValues)
}
