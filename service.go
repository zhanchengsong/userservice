package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zhanchengsong/userservice/route"
)

// This service handle user related request
// Incliding CRUD on user, friend/unfriend
func main() {
	// Loading extra env settings from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot not load env file")
	}
	http.Handle("/", route.Handlers())
	// Start service
	port := os.Getenv("SERVER_PORT")
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
