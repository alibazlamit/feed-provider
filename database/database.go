package database

import (
	reader "alibazlamit/feed-reader/feed-reader"
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()
var logger *log.Logger

func InitDatabase() *mongo.Collection {
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

	logger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	return collection
}

func GetAllArticlesFromDB() ([]reader.NewsArticleInformationMongoDB, error) {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Fatal("Error retrieving articles", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []reader.NewsArticleInformationMongoDB
	for cursor.Next(ctx) {
		var article reader.NewsArticleInformationMongoDB
		err := cursor.Decode(&article)
		if err != nil {
			logger.Fatal("Error decoding articles", err)
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func GetArticleByIDFromDB(id primitive.ObjectID) (*reader.NewsArticleInformationMongoDB, error) {
	filter := bson.M{"_id": id}
	var article reader.NewsArticleInformationMongoDB
	err := collection.FindOne(ctx, filter).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Fatal("Article not found", err)
			return nil, err
		}
		logger.Fatal("Error retrieving article", err)
		return nil, err
	}
	return &article, nil
}
