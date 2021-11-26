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
	routes.HandleFunc(createPath("register"), repository.RegisterUser).Methods("POST")
	routes.HandleFunc(createPath("login"), repository.ValidateUser).Methods("GET")
	routes.HandleFunc(createPath("users"), repository.ScanUsers).Methods("GET")
	routes.HandleFunc(createPath("users/{userId}"), repository.GetUser).Methods("GET")
	routes.HandleFunc(createPath("book/{bookId}"), apis.GetVolume)
	routes.HandleFunc(createPath("search"), apis.SearchBook)
	routes.Use(cors.MiddleWare)

	fmt.Printf("Server is running on port %s\n", "8083")
	http.ListenAndServe(":8083", routes)
}

func createPath(path string) string {
	return fmt.Sprintf("/%s/%s", API_BASE_PATH, path)
}
