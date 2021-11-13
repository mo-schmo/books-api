package apis

import (
	"booksApi/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

const GOOGLE_API_URL = "https://www.googleapis.com/books/v1"

func GetVolume(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookId := params["bookId"]
	url := fmt.Sprintf("%s/volumes/%s?key=%s", GOOGLE_API_URL, bookId, os.Getenv("GOOGLE_API_KEY"))
	res, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer res.Body.Close()

	book := entity.GoogleBook{}

	err = json.NewDecoder(res.Body).Decode(&book)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
	bookJson, err := json.Marshal(book)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bookJson)
}

func SearchBook(w http.ResponseWriter, r *http.Request) {
	query := Query{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Printf("Searching for book with query: %s\n", query.Query)
	var formattedQuery string
	if len(query.Query) > 0 {
		formattedQuery := strings.ReplaceAll(query.Query, " ", "+")
		fmt.Printf("Formatted query: %s\n", formattedQuery)
	}

	url := fmt.Sprintf("%s/volumes?q=%s&key=%s", GOOGLE_API_URL, formattedQuery, os.Getenv("GOOGLE_API_KEY"))
	fmt.Println(url)
}

type Query struct {
	Query string `json:"query"`
	Isbn  string `json:"isbn"`
}
