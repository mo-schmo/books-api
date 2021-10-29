package main

import (
	"booksApi/cors"
	"booksApi/repository"
	"fmt"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gorilla/mux"
)

func main() {
	routes := mux.NewRouter()
	routes.HandleFunc("/users", repository.ScanUsers).Methods("GET")

	fmt.Printf("Server is running on port %s\n", "8080")
	http.ListenAndServe(":8080", cors.MiddleWare(routes))
}
