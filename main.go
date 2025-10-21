package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
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
	db := connectPostgre()
	defer db.Close()

	insertSQL(db)
	selectSQL(db)

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

		saveOriginalUrl(shortenRequest.URL, shortID)

		log.Printf("добавил в мапу urlStore с ключем: %s и значением: %s", shortID, urlStore[shortID])
	} else {
		http.Error(w, "must be POST", http.StatusNotFound)
	}

}

func connectPostgre() *sql.DB {
	// Подключаемся к PostgreSQL
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=secret dbname=urlshortener sslmode=disable")
	if err != nil {
		log.Fatal("Не удалось открыть соединение с БД:", err)
	}

	// Проверяем живое подключение
	if err := db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	log.Println("Подключение к PostgreSQL установлено!")

	db.Exec(`DROP TABLE urls`) //удалить после дебага
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS urls (
        short_id TEXT PRIMARY KEY,
        original_url TEXT NOT NULL
    )
`)
	if err != nil {
		log.Fatal("Не удалось создать таблицу urls:", err)
	}
	log.Println("Таблица 'urls' готова")

	return db
}

func selectSQL(db *sql.DB) {
	rows, errQuery := db.Query("select * from urls")

	if errQuery != nil {
		log.Fatal("ЭТО ОШИБКА 2!", errQuery)
	}

	defer rows.Close()

	for rows.Next() {
		var shortID, originalURL string
		err := rows.Scan(&shortID, &originalURL)
		if err != nil {
			log.Fatal("Ошибка чтения строки:", err)
		}
		log.Printf("Найдено: %s → %s", shortID, originalURL)
	}

}

func insertSQL(db *sql.DB) {
	_, errExec := db.Exec(`INSERT INTO urls VALUES ('om0V4S', 'https://go.dev')`)
	if errExec != nil {
		log.Fatal("ЭТО ОШИБКА!", errExec)
	}
}

func redirectHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		id := strings.TrimPrefix(req.URL.Path, "/")
		originalURL, exists := getOriginalUrlById(id)
		if exists {
			http.Redirect(w, req, originalURL, http.StatusFound)
		} else {
			http.Error(w, "not found in urlStore", http.StatusNotFound)
		}
	} else {
		http.Error(w, "must be POST", http.StatusNotFound)
	}
}
func getOriginalUrlById(id string) (string, bool) {

	originalURL, exists := urlStore[id]
	return originalURL, exists
}

func saveOriginalUrl(url string, id string) {
	urlStore[id] = url
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
