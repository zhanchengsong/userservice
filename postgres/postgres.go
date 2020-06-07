package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"github.com/zhanchengsong/userservice/model"
	"os"
)

func ConnectDB() *gorm.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")
	databaseName := os.Getenv("PG_DBNAME")
	databaseHost := os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	// Get the url 
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)
	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		log.Fatal("DB connection failed")
		panic(err)
	}
	defer db.Close()
	// Migrate DDL
	db.AutoMigrate(
		&model.User
		)
	
}
