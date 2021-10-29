package main

import (
	"booksApi/cors"
	"booksApi/repository"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gorilla/mux"
)

func main() {
	routes := mux.NewRouter()
	routes.HandleFunc("/users", repository.ScanUsers).Methods("GET")
	http.ListenAndServe(":8080", cors.MiddleWare(routes))
}
