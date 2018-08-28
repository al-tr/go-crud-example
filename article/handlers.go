package article

import (
	"net/http"
	"strings"
)

const (
	Articles      = "/articles"
	ArticlesSlash = "/articles/"
)

func UrlArticle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetArticlesNotDeleted(w, r)
	case "PUT":
		PutArticle(w, r)
	case "DELETE":
		Clean(w, r)
	}
}

func UrlArticleSlash(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := strings.TrimPrefix(r.URL.Path, ArticlesSlash)
		if len(id) == 0 {
			GetArticlesNotDeleted(w, r)
			return
		}
		if id == "all" {
			GetArticlesAll(w, r)
			return
		}
		GetArticleById(w, r)
	case "PUT":
		PutArticle(w, r)
	case "DELETE":
		id := strings.TrimPrefix(r.URL.Path, ArticlesSlash)
		if len(id) == 0 {
			Clean(w, r)
			return
		}
		DeleteArticle(w, r)
	}
}
