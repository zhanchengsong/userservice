package route

import (
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"github.com/zhanchengsong/userservice/controller"
	_ "github.com/zhanchengsong/userservice/docs"
	"log"
	"net/http"
	"strings"
)

func Handlers() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(LoggingMiddleware)
	r.Use(CommonMiddleware)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	// These are the
	r.HandleFunc("/user", controller.CreateUser).Methods("POST")
	r.HandleFunc("/login", controller.Login).Methods("POST")
	r.HandleFunc("/user", controller.FindUser).Methods("GET")
	r.HandleFunc("/users", controller.FindUsersByPrefix).Methods("GET")
	r.HandleFunc("/follow", controller.CreateFollow).Methods("POST")
	r.HandleFunc("/followers", controller.GetFollowers).Methods("GET")
	r.HandleFunc("/followees", controller.GetFollowers).Methods("GET")
	return r
}

func CommonMiddleware(next http.Handler) http.Handler {
	// Middle ware takes a handler and returns a handler
	// The new handler is created with a http.HandlerFunc with a function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "swagger") {
			w.Header().Add("Content-Type", "application/json")
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request at", r.URL, r.Method)
		next.ServeHTTP(w, r)
	})
}
