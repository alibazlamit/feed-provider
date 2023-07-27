package main

import (
	"alibazlamit/feed-reader/database"
	reader "alibazlamit/feed-reader/feed-reader"
	"alibazlamit/feed-reader/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()
var logger *log.Logger
var articleRepository database.ArticleRepository

func main() {
	// Initialize MongoDB client or connection pool
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("news_feed").Collection("news")
	articleRepository = &database.MongoDBArticleRepository{
		Collection: collection,
		Logger:     logger,
	}

	r := reader.NewReader(articleRepository, logger, http.DefaultClient)

	err = r.RunCronFeedReader()
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
	responseObj := models.NewsArticlesResponse{
		Status: string(models.Failure),
	}
	articles, err := articleRepository.GetAllArticles()
	if err != nil {
		handleError(w, http.StatusBadRequest, "Error retrieving articles", err)
	}
	responseObj = models.NewsArticlesResponse{
		Data:   articles,
		Status: string(models.Success),
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

	article, err := articleRepository.GetArticleByID(objectID)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Error retrieving article", err)
	}

	responseObj := models.NewsArticleResponse{
		Status: string(models.Success),
		Data:   *article,
	}
	handleSuccess(w, http.StatusOK, responseObj)
}

func handleError(w http.ResponseWriter, statusCode int, message string, err error) {
	logger.Printf("Error: %v", err)
	responseObj := models.NewsArticlesResponse{
		Error:  message,
		Status: string(models.Failure),
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
