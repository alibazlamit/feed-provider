package main

import (
	"alibazlamit/feed-provider/database"
	"alibazlamit/feed-provider/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllArticles(t *testing.T) {
	mockRepo := database.NewMockArticleRepository()
	articleRepository = mockRepo

	article1 := models.NewsArticleInformationMongoDB{
		Title:    "Article 1",
		BodyText: "Content 1",
	}
	article2 := models.NewsArticleInformationMongoDB{
		Title:    "Article 2",
		BodyText: "Content 2",
	}
	mockRepo.Articles = append(mockRepo.Articles, article1, article2)

	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		getAllArticles(w, r)
	}).Methods("GET")

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	var responseObj models.NewsArticlesResponse
	err = json.Unmarshal(rr.Body.Bytes(), &responseObj)
	if err != nil {
		t.Fatal(err)
	}

	if len(responseObj.Data) != 2 {
		t.Errorf("Expected 2 articles, but got %d", len(responseObj.Data))
	}
	if responseObj.Data[0].Title != article1.Title {
		t.Errorf("Expected article title %s, but got %s", article1.Title, responseObj.Data[0].Title)
	}
	if responseObj.Data[1].Title != article2.Title {
		t.Errorf("Expected article title %s, but got %s", article2.Title, responseObj.Data[1].Title)
	}
}

func TestGetArticleByID(t *testing.T) {
	router := mux.NewRouter()
	id := primitive.NewObjectID()

	mockRepo := database.NewMockArticleRepository()
	mockRepo.Articles = []models.NewsArticleInformationMongoDB{
		models.NewsArticleInformationMongoDB{
			ID:       id,
			Title:    "Test Article",
			BodyText: "This is a test article.",
		},
	}
	articleRepository = mockRepo

	req, err := http.NewRequest("GET", fmt.Sprintf("/articles/%v", id.Hex()), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getArticleByID(w, r)
	})

	router.Handle("/articles/{id}", handler).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var responseObj models.NewsArticleResponse
	err = json.Unmarshal(rr.Body.Bytes(), &responseObj)
	if err != nil {
		t.Errorf("error decoding response body: %v", err)
	}

	expected := &models.NewsArticleInformationMongoDB{
		ID:       id,
		Title:    "Test Article",
		BodyText: "This is a test article.",
	}
	if responseObj.Data.ID != expected.ID || responseObj.Data.Title != expected.Title || responseObj.Data.BodyText != expected.BodyText {
		t.Errorf("handler returned unexpected body: got %v want %v", responseObj, expected)
	}
}
