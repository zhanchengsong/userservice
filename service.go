package main
import (
	"github.com/zhanchengsong/userservice/postgres"
)
// This service handle user related request
// Incliding CRUD on user, friend/unfriend
func main() {
	postgres.ConnectDB()
}
