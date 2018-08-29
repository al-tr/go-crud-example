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

	log.Print("Get article by id")

	id := strings.TrimPrefix(r.URL.Path, articlesSlashUrl)
	if len(id) == 0 {
		getArticlesNotDeletedService(w, r)
		return
	}

	log.Print("Get article by id: '", id, "'")

	articleById, err := getArticleById(id)
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
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

	articleFromDatabase, err := getArticleById(*article.Uuid)
	if err != nil {
		// TODO: create
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}
	// TODO update

	log.Print("Found in db: ", articleFromDatabase)

	var errors []string
	if stringNilOrEmpty(articleFromDatabase.Text) {
		errors = append(errors, "text must not be empty")
	}
	if stringNilOrEmpty(articleFromDatabase.Title) {
		errors = append(errors, "title must not be empty")
	}
	if stringNilOrEmpty(articleFromDatabase.Publisher) {
		errors = append(errors, "publisher must not be empty")
	}

	if len(errors) > 0 {
		createErrorResponse(w, 400, errors)
		return
	}

	var articleToInsertIntoDb Article
	id := uuid.New().String()
	articleToInsertIntoDb.Uuid = &id
	articleToInsertIntoDb.Title = article.Title
	articleToInsertIntoDb.Text = article.Text
	articleToInsertIntoDb.Publisher = user.Email
	formattedTime := nowUtc()
	articleToInsertIntoDb.DatePublished = &formattedTime
	if article.IsDeleted != nil {
		articleToInsertIntoDb.IsDeleted = article.IsDeleted
	}

	insertedArticle, err := putArticle(articleToInsertIntoDb)
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	articleJson, err := json.Marshal(insertedArticle)
	if err != nil {
		createErrorResponse(w, 500, []string{err.Error()})
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write(articleJson)
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
			if articleFromDb.Publisher == nil || len(*articleFromDb.Publisher) == 0 ||
				articleFromDb.Title == nil || len(*articleFromDb.Title) == 0 ||
				articleFromDb.Text == nil || len(*articleFromDb.Text) == 0 ||
				articleFromDb.Publisher == nil || len(*articleFromDb.Publisher) == 0 ||
				articleFromDb.DatePublished == nil || len(*articleFromDb.DatePublished) == 0 {

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
	updateArticles, ers := bulkPutArticles(articlesToUpdate)

	var response MultipleResponse
	for i := 0; i < len(updateArticles); i++ {
		var data DataResponse
		data.Id = updateArticles[i].Uuid
		errors := []string{ers[i].Error()}
		data.Errors = &errors

		responseTmp := append(*response.Responses, data)
		response.Responses = &responseTmp
	}
	createDataStringResponse(w, 200, "Docs deleted during cleanService: "+strconv.Itoa(numberOfDeletedDocs))
}

func deleteArticleService(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, articlesUrl)
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

	deleteArticleById(id, *user)
	createDataStringResponse(w, 200, "everything's good for id")
}
