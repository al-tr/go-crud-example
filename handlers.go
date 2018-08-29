package main

import (
	"net/http"
	"strings"
)

const (
	articlesUrl      = "/articles"
	articlesSlashUrl = "/articles/"
)

func urlArticle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getArticlesNotDeletedService(w, r)
	case "POST":
		postArticleService(w, r)
	case "DELETE":
		cleanService(w, r)
	}
}

func urlArticleSlash(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := strings.TrimPrefix(r.URL.Path, articlesSlashUrl)
		if len(id) == 0 {
			getArticlesNotDeletedService(w, r)
			return
		}
		if id == "all" {
			getArticlesAllService(w, r)
			return
		}
		getArticleByIdService(w, r)
	case "POST":
		postArticleService(w, r)
	case "PUT":
		putArticleService(w, r)
	case "DELETE":
		id := strings.TrimPrefix(r.URL.Path, articlesSlashUrl)
		if len(id) == 0 {
			cleanService(w, r)
			return
		}
		deleteArticleService(w, r)
	}
}
