package main

import (
	reader "alibazlamit/feed-reader/feed-reader"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

func main() {
	err := reader.RunCronFeedReader(collection)
	if err != nil {
		fmt.Println("Error running feed reader job:", err)
	}

	// Create a new router instance
	router := mux.NewRouter()

	// Define the HTTP route handlers
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
	router.HandleFunc("/articles", getAllArticles).Methods("GET")
	router.HandleFunc("/articles/{id}", getArticleByID).Methods("GET")

	// Start the HTTP server on port 8080
	fmt.Println("Server listening on http://localhost:8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error:", err)
	}

}

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("news_feed").Collection("news")
}

// GetAllArticles returns all articles from the MongoDB database in JSON format
func getAllArticles(w http.ResponseWriter, r *http.Request) {
	responseObj := reader.NewsArticlesResponse{
		Status: string(reader.Failure),
	}
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Error retrieving articles", err)
		return
	}
	defer cursor.Close(ctx)

	var articles []reader.NewsArticleInformationMongoDB
	for cursor.Next(ctx) {
		var article reader.NewsArticleInformationMongoDB
		err := cursor.Decode(&article)
		if err != nil {
			handleError(w, http.StatusInternalServerError, "Error decoding articles", err)
			return
		}
		articles = append(articles, article)
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

	filter := bson.M{"_id": objectID}
	var article reader.NewsArticleInformationMongoDB
	err = collection.FindOne(ctx, filter).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			handleError(w, http.StatusNotFound, "Article not found", err)
			return
		}
		handleError(w, http.StatusInternalServerError, "Error retrieving article", err)
		return
	}

	responseObj := reader.NewsArticleResponse{
		Status: string(reader.Success),
		Data:   article,
	}
	handleSuccess(w, http.StatusOK, responseObj)
}

func handleError(w http.ResponseWriter, statusCode int, message string, err error) {
	// Log the error (you can use a logger library for proper error logging)
	fmt.Println("Error:", err)

	// Prepare the response object
	responseObj := reader.NewsArticlesResponse{
		Error:  message,
		Status: string(reader.Failure),
	}

	// Encode the response object and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(responseObj); err != nil {
		// If encoding the error response fails, fallback to simple text response
		w.Write([]byte(message))
	}
}

func handleSuccess(w http.ResponseWriter, statusCode int, responseObj interface{}) {
	// Encode the response object and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(responseObj); err != nil {
		// Handle the error and send response
		handleError(w, http.StatusInternalServerError, "Error encoding response", err)
		return
	}
}
