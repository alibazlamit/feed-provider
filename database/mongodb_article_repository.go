package database

import (
	"alibazlamit/feed-provider/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	NEWS_ARTICLE_KEY = "NewsArticleID"
)

var ctx = context.TODO()

type MongoDBArticleRepository struct {
	Collection *mongo.Collection
	Logger     *log.Logger
}

func (r *MongoDBArticleRepository) GetArticleByID(id primitive.ObjectID) (*models.NewsArticleInformationMongoDB, error) {
	filter := bson.M{"_id": id}
	var article models.NewsArticleInformationMongoDB
	err := r.Collection.FindOne(ctx, filter).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.Logger.Fatal("Article not found", err)
			return nil, err
		}
		r.Logger.Fatal("Error retrieving article", err)
		return nil, err
	}
	return &article, nil
}

func (r *MongoDBArticleRepository) GetAllArticles() ([]models.NewsArticleInformationMongoDB, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		r.Logger.Fatal("Error retrieving articles", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []models.NewsArticleInformationMongoDB
	for cursor.Next(ctx) {
		var article models.NewsArticleInformationMongoDB
		err := cursor.Decode(&article)
		if err != nil {
			r.Logger.Fatal("Error decoding articles", err)
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (r *MongoDBArticleRepository) AddOrUpdateArticle(articleID int, articleXml *models.NewsArticleInformationXML) error {
	opts := options.Replace().SetUpsert(true)
	filter := bson.D{{Key: NEWS_ARTICLE_KEY, Value: articleID}}
	_, err := r.Collection.ReplaceOne(context.TODO(), filter, models.ConvertToMongoDB(articleXml), opts)
	if err != nil {
		r.Logger.Printf("Error saving article with id:%d and error: %v\n", articleID, err)
		return err
	}
	return nil
}
