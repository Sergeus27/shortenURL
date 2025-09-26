package main

import (
	"io"
	"log"
	"net/http"
)

func main() {

	shortenHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			io.WriteString(w, "OK\n")
		} else {
			http.Error(w, "must be POST", http.StatusNotFound)
		}
	}

	http.HandleFunc("/api/shorten", shortenHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
