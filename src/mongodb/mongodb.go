package mongodb

import (
	"log"

	"github.com/go-bongo/bongo"
)

// CreateMongoConnection returns a mongo db connection upon success
func CreateMongoConnection(connectionString string, dbName string) bongo.Connection {
	config := &bongo.Config{
		ConnectionString: connectionString,
		Database:         dbName,
	}
	connection, err := bongo.Connect(config)
	if err != nil {
		log.Fatal(err)
	}
	return connection
}
