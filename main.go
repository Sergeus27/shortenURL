package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
)

type ShortenRequest struct {
	URL string `json:url`
}

func main() {
	log.Printf("ID:%s", generateShortID())
	shortenHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			var shortenRequest ShortenRequest
			//объявляется error создается декодер и используется метод decode из этого декодера в который передаем ссылку на переменную для сохранения декодированного Body
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(&shortenRequest); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			log.Printf("URL from shorten handle - %s", shortenRequest.URL)
		} else {
			http.Error(w, "must be POST", http.StatusNotFound)
		}

	}
	//почитать как работает http сервер(жизненный цикл)
	http.HandleFunc("/api/shorten", shortenHandler)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
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
