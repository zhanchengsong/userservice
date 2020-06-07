package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

// CreateClient is a function that initialize the the firebase client
func CreateClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := "mytwitter-279300"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}
