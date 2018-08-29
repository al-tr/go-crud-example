package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getArticlesNotDeletedService(w http.ResponseWriter, r *http.Request) {
	user, e := authenticateRequest(r)
	if e != nil {
		createErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		createErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	log.Print("Get articles not deleted")

	articlesFromDatabase, err := getArticlesNotDeleted()
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	articles := articlesFromDatabase

	log.Print("Got ", len(articles), " articles")

	if len(articles) == 0 {
		articles = make([]Article, 0) // workaround to get [] instead of null after json.Marshal
	}
	articleJson, err := json.Marshal(articles)
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write(articleJson)
}

func getArticlesAllService(w http.ResponseWriter, r *http.Request) {
	user, e := authenticateRequest(r)
	if e != nil {
		createErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		createErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	log.Print("Get articles all")

	articles, err := getAllArticles()
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	log.Print("Got ", len(articles), " articles")

	if len(articles) == 0 {
		articles = make([]Article, 0) // workaround to get [] instead of null after json.Marshal
	}
	articleJson, err := json.Marshal(articles)
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write(articleJson)
}

func getArticleByIdService(w http.ResponseWriter, r *http.Request) {
	user, e := authenticateRequest(r)
	if e != nil {
		createErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		createErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	id := strings.TrimPrefix(r.URL.Path, articlesSlashUrl)
	if len(id) == 0 {
		getArticlesNotDeletedService(w, r)
		return
	}

	log.Print("Get article by id: '", id, "'")

	articleById, err := getArticleById(id)
	if err != nil {
		createErrorResponse(w, 400, []string{err.Error()})
		return
	}

	articleJson, err := json.Marshal(articleById)
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write(articleJson)
}

func postArticleService(w http.ResponseWriter, r *http.Request) {
	user, e := authenticateRequest(r)
	if e != nil {
		createErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		createErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	idFromUrl := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, articlesUrl), "/")
	if len(idFromUrl) > 0 {
		createErrorResponse(w, 400, []string{"post with idFromUrl is not supported, use put"})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		createErrorResponse(w, 400, []string{err.Error()})
		return
	}

	log.Print("Post article '", body, "'")

	var article *Article
	err = json.Unmarshal(body, &article)
	if err != nil {
		createErrorResponse(w, 400, []string{err.Error()})
		return
	}

	errors := validateArticle(article)
	if len(errors) > 0 {
		createErrorResponse(w, 400, errors)
		return
	}

	now := nowUtc()
	id := uuid.New().String()
	var articleToInsertIntoDb = Article{Uuid: &id, Title: article.Title, Text: article.Text, Publisher: user.Email, DatePublished: &now, IsDeleted: article.IsDeleted}

	createdArticle, err := putArticle(articleToInsertIntoDb)
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	articleJson, err := json.Marshal(createdArticle)
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write(articleJson)
}

func putArticleService(w http.ResponseWriter, r *http.Request) {
	user, e := authenticateRequest(r)
	if e != nil {
		createErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		createErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	log.Print("Put article")

	id := strings.TrimPrefix(r.URL.Path, articlesSlashUrl)
	if len(id) == 0 {
		createErrorResponse(w, 400, []string{"provide id in path"})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		createErrorResponse(w, 400, []string{err.Error()})
		return
	}

	log.Print("Put article '", body, "'")

	var article *Article
	err = json.Unmarshal(body, &article)
	if err != nil {
		createErrorResponse(w, 400, []string{err.Error()})
		return
	}

	articleFromDatabase, err := getArticleById(id)
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	log.Print("Found in db: ", articleFromDatabase)

	if article.Title != nil {
		articleFromDatabase.Title = article.Title
	}
	if article.Text != nil {
		articleFromDatabase.Text = article.Text
	}
	if article.IsDeleted != nil {
		articleFromDatabase.IsDeleted = article.IsDeleted
	}
	articleFromDatabase.Editor = user.Email
	formattedTime := nowUtc()
	articleFromDatabase.DateUpdated = &formattedTime


	updatedArticle, err := putArticle(*articleFromDatabase)
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	articleJson, err := json.Marshal(updatedArticle)
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write(articleJson)
}

func validateArticle(article *Article) []string {
	var errors []string
	if stringNilOrEmpty(article.Text) {
		errors = append(errors, "text must not be empty")
	}
	if stringNilOrEmpty(article.Title) {
		errors = append(errors, "title must not be empty")
	}
	return errors
}

func cleanService(w http.ResponseWriter, r *http.Request) {
	user, e := authenticateRequest(r)
	if e != nil {
		createErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		createErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	articles, err := getArticlesNotDeleted()
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	var numberOfDeletedDocs = 0
	var articlesToUpdate []Article
	for i := 0; i < len(articles); i++ {
		articleFromDb := articles[i]

		if articleFromDb.IsDeleted == nil || *articleFromDb.IsDeleted == false {
			if stringNilOrEmpty(articleFromDb.Publisher) ||
				stringNilOrEmpty(articleFromDb.Title) ||
				stringNilOrEmpty(articleFromDb.Text) ||
				stringNilOrEmpty(articleFromDb.Publisher) ||
				stringNilOrEmpty(articleFromDb.DatePublished) {

				numberOfDeletedDocs++

				booleanTrue := true
				articleFromDb.IsDeleted = &booleanTrue
				articleFromDb.Editor = user.Email
				timeNow := nowUtc()
				articleFromDb.DateUpdated = &timeNow

				articlesToUpdate = append(articlesToUpdate, articleFromDb)
			}
		}
	}
	updatedArticles, ers := bulkPutArticles(articlesToUpdate)

	var response MultipleResponse
	for i := 0; i < len(updatedArticles); i++ {
		var data DataResponse
		data.Id = updatedArticles[i].Uuid
		errors := []string{ers[i].Error()}
		data.Errors = &errors

		responseTmp := append(*response.Responses, data)
		response.Responses = &responseTmp
	}
	createDataStringResponse(w, 200, "Docs deleted during cleanService: "+strconv.Itoa(numberOfDeletedDocs))
}

func deleteArticleService(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, articlesSlashUrl)
	if len(id) == 0 {
		cleanService(w, r)
		return
	}

	user, e := authenticateRequest(r)
	if e != nil {
		createErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		createErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	article, err := getArticleById(id)
	if err != nil {
		createErrorResponse(w, 400, []string{err.Error()})
		return
	}

	deleted := true
	article.IsDeleted = &deleted
	article.Editor = user.Email
	now := nowUtc()
	article.DateUpdated = &now

	_, err = putArticle(*article)
	if err != nil {
		createErrorResponse(w, 400, []string{err.Error()})
		return
	}

	createDataStringResponse(w, 200, "article with id '"+id+"' is marked as deleted")
}
