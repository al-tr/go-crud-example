package main

import (
	"crud/article"
	"log"
	"net/http"
	"os"
)

func main() {
	article.InitDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		defaultPort := "8080"
		log.Print("$PORT is not set, using ", defaultPort)
		port = defaultPort
	}

	http.HandleFunc(article.Articles, article.UrlArticle)
	http.HandleFunc(article.ArticlesSlash, article.UrlArticleSlash)
	http.ListenAndServe(":"+port, nil)
}
