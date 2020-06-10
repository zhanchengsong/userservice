package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zhanchengsong/userservice/model"
	"log"
)

func ConnectDB(username string, password string, databaseName string, databaseHost string ) *gorm.DB {
	//err := godotenv.Load("../.env")
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	//username := os.Getenv("PG_USERNAME")
	//password := os.Getenv("PG_PASSWORD")
	//databaseName := os.Getenv("PG_DBNAME")
	//databaseHost := os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	// Get the url 
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)
	log.Print("Connecting using uri", dbURI)
	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		log.Fatal("DB connection failed")
		panic(err)
	}
	defer db.Close()
	// Migrate DDL
	db.AutoMigrate(
		&model.User{})
	log.Print("Succesfully connected to db")
	return db
}
