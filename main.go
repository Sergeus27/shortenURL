package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ShortenRequest struct {
	URL string `json:url`
}

func main() {

	shortenHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			var shortenRequest ShortenRequest
			if err := json.NewDecoder(req.Body).Decode(&shortenRequest); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			log.Printf("URL from shorten handle - %s", shortenRequest.URL)
		} else {
			http.Error(w, "must be POST", http.StatusNotFound)
		}

	}

	http.HandleFunc("/api/shorten", shortenHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
