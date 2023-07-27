package database

import (
	"alibazlamit/feed-provider/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockArticleRepository struct {
	Articles []models.NewsArticleInformationMongoDB
}

func NewMockArticleRepository() *MockArticleRepository {
	return &MockArticleRepository{
		Articles: []models.NewsArticleInformationMongoDB{},
	}
}

func (r *MockArticleRepository) GetAllArticles() ([]models.NewsArticleInformationMongoDB, error) {
	return r.Articles, nil
}

func (r *MockArticleRepository) GetArticleByID(id primitive.ObjectID) (*models.NewsArticleInformationMongoDB, error) {
	for _, article := range r.Articles {
		if article.ID == id {
			return &article, nil
		}
	}
	return nil, nil // Return nil if article not found
}

func (r *MockArticleRepository) AddOrUpdateArticle(id int, article *models.NewsArticleInformationXML) error {
	newsArticle := models.NewsArticleInformationMongoDB{
		ID:             primitive.NewObjectID(),
		Title:          article.NewsArticle.Title,
		BodyText:       article.NewsArticle.BodyText,
		ClubName:       article.ClubName,
		ClubWebsiteURL: article.ClubWebsiteURL,
		ArticleURL:     article.NewsArticle.ArticleURL,
		NewsArticleID:  article.NewsArticle.NewsArticleID,
		Taxonomies:     article.NewsArticle.Taxonomies,
	}

	r.Articles = append(r.Articles, newsArticle)
	return nil
}
