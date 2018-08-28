package article

import (
	"crud/auth"
	"crud/util"
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetArticlesNotDeleted(w http.ResponseWriter, r *http.Request) {
	user, e := auth.AuthenticateRequest(r)
	if e != nil {
		util.CreateErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		util.CreateErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	log.Print("Get articles not deleted")

	articles := getAllArticlesNotDeleted()

	log.Print("Got ", len(articles), " articles")

	articleJson, err := json.Marshal(articles)
	util.Panicerr(err)

	w.Header().Add("Content-type", "application/json")
	w.Write(articleJson)
}

func GetArticlesAll(w http.ResponseWriter, r *http.Request) {
	user, e := auth.AuthenticateRequest(r)
	if e != nil {
		util.CreateErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		util.CreateErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	log.Print("Get articles all")

	articles := getAllArticles()

	log.Print("Got ", len(articles), " articles")

	articleJson, err := json.Marshal(articles)
	util.Panicerr(err)

	w.Header().Add("Content-type", "application/json")
	w.Write(articleJson)
}

func GetArticleById(w http.ResponseWriter, r *http.Request) {
	user, e := auth.AuthenticateRequest(r)
	if e != nil {
		util.CreateErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		util.CreateErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	log.Print("Get article by id")

	id := strings.TrimPrefix(r.URL.Path, ArticlesSlash)
	if len(id) == 0 {
		GetArticlesNotDeleted(w, r)
		return
	}

	log.Print("Get article by id: '", id, "'")

	articleById := getArticleById(id)
	articleJson, err := json.Marshal(articleById)
	util.Panicerr(err)

	w.Header().Add("Content-type", "application/json")
	w.Write(articleJson)
}

func PutArticle(w http.ResponseWriter, r *http.Request) {
	user, e := auth.AuthenticateRequest(r)
	if e != nil {
		util.CreateErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		util.CreateErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	log.Print("Put article")

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		util.CreateErrorResponse(w, 400, []string{err.Error()})
		return
	}

	log.Print("Put article '", body, "'")

	var article *Article
	err = json.Unmarshal(body, &article)
	if err != nil {
		util.CreateErrorResponse(w, 400, []string{err.Error()})
		return
	}

	articleFromDatabase := getArticleById(*article.Uuid)
	log.Print("Found in db: ", articleFromDatabase)

	var errors []string
	if util.StringNilOrEmpty(article.Text) {
		errors = append(errors, "text must not be empty")
	}
	if util.StringNilOrEmpty(article.Title) {
		errors = append(errors, "title must not be empty")
	}
	if util.StringNilOrEmpty(article.Publisher) {
		errors = append(errors, "publisher must not be empty")
	}

	if len(errors) > 0 {
		util.CreateErrorResponse(w, 400, errors)
		return
	}

	var articleToInsertIntoDb Article
	id := uuid.New().String()
	articleToInsertIntoDb.Uuid = &id
	articleToInsertIntoDb.Title = article.Title
	articleToInsertIntoDb.Text = article.Text
	articleToInsertIntoDb.Publisher = user.Email
	formattedTime := util.NowUtc()
	articleToInsertIntoDb.DatePublished = &formattedTime
	if article.IsDeleted != nil {
		articleToInsertIntoDb.IsDeleted = article.IsDeleted
	}

	insertedArticle := updateArticle(articleToInsertIntoDb)

	articleJson, err := json.Marshal(insertedArticle)
	if err != nil {
		util.CreateErrorResponse(w, 500, []string{err.Error()})
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.Write(articleJson)
}

func Clean(w http.ResponseWriter, r *http.Request) {
	user, e := auth.AuthenticateRequest(r)
	if e != nil {
		util.CreateErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		util.CreateErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	articles := getAllArticlesNotDeleted()

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
				timeNow := util.NowUtc()
				articleFromDb.DateUpdated = &timeNow

				articlesToUpdate = append(articlesToUpdate, articleFromDb)
			}
		}
	}
	bulkUpdateArticles(articlesToUpdate)
	util.CreateDataStringResponse(w, 200, "Docs deleted during Clean: "+strconv.Itoa(numberOfDeletedDocs))
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, Articles)
	if len(id) == 0 {
		Clean(w, r)
		return
	}

	user, e := auth.AuthenticateRequest(r)
	if e != nil {
		util.CreateErrorResponse(w, 401, []string{"not authenticated, provide email in header 'Authorization:Bearer email@example.com'"})
		return
	}

	if user.Email == nil {
		util.CreateErrorResponse(w, 403, []string{"not authorized"})
		return
	}

	deleteArticleById(id)
	util.CreateDataStringResponse(w, 200, "everything's good for id")
}