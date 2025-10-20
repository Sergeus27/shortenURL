package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

var urlStore = make(map[string]string)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortID string `json:"short_id"`
}

func main() {
	// log.Printf("ID:%s", generateShortID())
	//почитать как работает http сервер(жизненный цикл)
	http.HandleFunc("/api/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)

}

func shortenHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		var shortenRequest ShortenRequest
		//объявляется error создается декодер и используется метод decode из этого декодера в который передаем ссылку на переменную для сохранения декодированного Body
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&shortenRequest); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		log.Printf("URL from shorten handle - %s", shortenRequest.URL)

		shortID := generateShortID()
		response := ShortenResponse{ShortID: shortID}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

		urlStore[shortID] = shortenRequest.URL
		log.Printf("добавил в мапу urlStore с ключем: %s и значением: %s", shortID, urlStore[shortID])
	} else {
		http.Error(w, "must be POST", http.StatusNotFound)
	}

}

func redirectHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		id := strings.TrimPrefix(req.URL.Path, "/")
		originalURL, exists := urlStore[id]
		if exists {
			http.Redirect(w, req, originalURL, http.StatusFound)
		} else {
			http.Error(w, "not found in urlStore", http.StatusNotFound)
		}
	} else {
		http.Error(w, "must be POST", http.StatusNotFound)
	}
}

func generateShortID() string {
	symbols := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 6
	maxLen := len(symbols)
	shortID := ""
	for i := 0; i < length; i++ {
		shortID = shortID + string(symbols[rand.Intn(maxLen)])
	}
	return shortID
}
