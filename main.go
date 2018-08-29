package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "my.db"
		log.Print("$DB_NAME is not set, using ", dbName)
	}
	initDatabase(dbName, true)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Print("$PORT is not set, using ", port)
	}

	http.HandleFunc(articlesUrl, urlArticle)
	http.HandleFunc(articlesSlashUrl, urlArticleSlash)
	http.ListenAndServe(":"+port, nil)
}
