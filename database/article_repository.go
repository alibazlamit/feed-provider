package database

import (
	"alibazlamit/feed-reader/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArticleRepository interface {
	GetArticleByID(id primitive.ObjectID) (*models.NewsArticleInformationMongoDB, error)
	GetAllArticles() ([]models.NewsArticleInformationMongoDB, error)
	AddOrUpdateArticle(articleID int, articleXml *models.NewsArticleInformationXML) error
}
