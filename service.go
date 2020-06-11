package main

import (
	"github.com/joho/godotenv"
	"github.com/zhanchengsong/userservice/route"
	"log"
	"net/http"
)
// This service handle user related request
// Incliding CRUD on user, friend/unfriend
func main() {
	// Loading extra env settings from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot not load env file")
	}
	http.Handle("/", route.Handlers())
	// Start service
	port := "3001"
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
