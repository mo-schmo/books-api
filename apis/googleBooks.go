package apis

import (
	"booksApi/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
