package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	dbName = "test.db"
)

func setup() {
	deleteTestDatabase(dbName)
	initDatabase(dbName, false)
}

func deleteTestDatabase(dbName string) {
	file, _ := os.Open(dbName + ".lock")
	if file != nil {
		file.Close()
		err := os.Remove(dbName + ".lock")
		if err != nil {
			log.Print(err.Error())
		}
	}
	file, _ = os.Open(dbName)
	if file != nil {
		file.Close()
		err := os.Remove(dbName)
		if err != nil {
			log.Print(err.Error())
		}
	}
}

func tearDown() {
	err := closeConnection()
	if err != nil {
		log.Print("Error while closing " + err.Error())
	}
	deleteTestDatabase(dbName)
}

func TestGetArticles(t *testing.T) {
	setup()

	t.Run("returns nothing", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, articlesUrl, nil)
		request.Header.Add("Authorization", "Bearer tester@rago.com")
		response := httptest.NewRecorder()

		getArticlesAllService(response, request)

		got := response.Body.String()
		want := "[]"

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	})

	tearDown()
}

func TestPostArticles(t *testing.T) {
	setup()

	t.Run("returns nothing", func(t *testing.T) {
		jsonBody := []byte(`{"title":"Finally, a new title","text":"Some text for real fun"}`)
		request, _ := http.NewRequest(http.MethodPost, articlesUrl, bytes.NewBuffer(jsonBody))
		request.Header.Add("Authorization", "Bearer tester@rago.com")
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		postArticleService(response, request)

		got := response.Body.String()
		var article *Article
		err := json.Unmarshal([]byte(got), &article)
		if err != nil {
			t.Error(err.Error())
		}

		assertEquals(t, *article.Title, "Finally, a new title")
		assertEquals(t, *article.Text, "Some text for real fun")
		assertEquals(t, *article.Publisher, "tester@rago.com")
		if stringNilOrEmpty(article.Uuid) {
			t.Error("Uuid is nil")
		}
		if stringNilOrEmpty(article.DatePublished) {
			t.Error("DatePublished is nil")
		}
	})

	tearDown()
}

func assertEquals(t *testing.T, got string, want string) {
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
