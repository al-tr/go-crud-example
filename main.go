package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	InitDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		defaultPort := "8080"
		log.Print("$PORT is not set, using ", defaultPort)
		port = defaultPort
	}

	http.HandleFunc(articlesUrl, urlArticle)
	http.HandleFunc(articlesSlashUrl, urlArticleSlash)
	http.ListenAndServe(":"+port, nil)
}
