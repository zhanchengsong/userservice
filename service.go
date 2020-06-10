package main
import (
	"github.com/joho/godotenv"
	"github.com/zhanchengsong/userservice/postgres"
	"log"
	"os"
)
// This service handle user related request
// Incliding CRUD on user, friend/unfriend
func main() {
	// Loading extra env settings from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot not load env file")
	}
	username := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")
	databaseName := os.Getenv("PG_DBNAME")
	databaseHost := os.Getenv("DB_HOST")
	postgres.ConnectDB(username, password, databaseName, databaseHost)
}
