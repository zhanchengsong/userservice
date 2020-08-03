package main

import (
	"log"
	"net/http"
	"os"
	"github.com/rs/cors"
	"github.com/joho/godotenv"
	"github.com/zhanchengsong/userservice/route"
)

// This service handle user related request
// Incliding CRUD on user, friend/unfriend
func main() {
	// Loading extra env settings from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("serverice.go: Cannot not load env file")
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowCredentials: true,
		});

	http.Handle("/", c.Handler ( router ) );
	// Start service
	port := os.Getenv("SERVER_PORT")
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
