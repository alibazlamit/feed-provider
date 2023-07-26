package main

import (
	"alibazlamit/feed-reader/database"
	reader "alibazlamit/feed-reader/feed-reader"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ctx = context.TODO()
var logger *log.Logger

func main() {
	r := reader.NewReader(database.InitDatabase(), logger)

	err := r.RunCronFeedReader()
	if err != nil {
		logger.Fatalf("Error running cron feed reader: %v", err)
	}

	router := mux.NewRouter()

	// API endpoints
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "PONG")
	})
	router.HandleFunc("/articles", getAllArticles).Methods("GET")
	router.HandleFunc("/articles/{id}", getArticleByID).Methods("GET")

	// Start the HTTP server on port 8080
	fmt.Println("Server listening on http://localhost:8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Fatalf("Error: %v", err)
	}

}

// GetAllArticles returns all articles from the MongoDB database in JSON format
func getAllArticles(w http.ResponseWriter, r *http.Request) {
	responseObj := reader.NewsArticlesResponse{
		Status: string(reader.Failure),
	}
	articles, err := database.GetAllArticlesFromDB()
	if err != nil {
		handleError(w, http.StatusBadRequest, "Error retrieving articles", err)
	}
	responseObj = reader.NewsArticlesResponse{
		Data:   articles,
		Status: string(reader.Success),
	}

	handleSuccess(w, http.StatusOK, responseObj)
}

// GetArticleByID returns the article with the specified ID from the MongoDB database
func getArticleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid article ID", err)
		return
	}

	article, err := database.GetArticleByIDFromDB(objectID)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Error retrieving article", err)
	}

	responseObj := reader.NewsArticleResponse{
		Status: string(reader.Success),
		Data:   *article,
	}
	handleSuccess(w, http.StatusOK, responseObj)
}

func handleError(w http.ResponseWriter, statusCode int, message string, err error) {
	logger.Printf("Error: %v", err)
	responseObj := reader.NewsArticlesResponse{
		Error:  message,
		Status: string(reader.Failure),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(responseObj); err != nil {
		w.Write([]byte(message))
	}
}

func handleSuccess(w http.ResponseWriter, statusCode int, responseObj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(responseObj); err != nil {
		handleError(w, http.StatusInternalServerError, "Error encoding response", err)
		return
	}
}
