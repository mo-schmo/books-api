package main

import (
	"booksApi/apis"
	"booksApi/cors"
	"booksApi/repository"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

const API_BASE_PATH = "api"

func main() {
	routes := mux.NewRouter()
	routes.HandleFunc(createPath("users"), repository.ScanUsers).Methods("GET")
	routes.HandleFunc("/api/book/{bookId}", apis.GetVolume)

	fmt.Printf("Server is running on port %s\n", "8083")
	http.ListenAndServe(":8083", cors.MiddleWare(routes))
}

func createPath(path string) string {
	return fmt.Sprintf("/%s/%s", API_BASE_PATH, path)
}
