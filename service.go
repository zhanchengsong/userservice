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
// Including CRUD on user, friend/unfriend
// @title Userservice API
// @version 1.0
// @description This microservice supports basic user related operations

// @contact.name Zhancheng Song, Gordon Lee

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost/api/v1/users
// @BasePath /
// @query.collection.format multi

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

	http.Handle("/", c.Handler ( route.Handlers() ) );
	// Start service
	port := os.Getenv("SERVER_PORT")
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
