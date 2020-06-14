package route
import (
	"github.com/gorilla/mux"
	"github.com/zhanchengsong/userservice/controllers"
	"net/http"

)
func Handlers() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)
	// These are the
	r.HandleFunc("/user", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/token", controllers.Login).Methods("POST")
	return r
}

func CommonMiddleware(next http.Handler) http.Handler {
	// Middle ware takes a handler and returns a handler
	// The new handler is created with a http.HandlerFunc with a function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}