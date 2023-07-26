package reader

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	NEWS_ARTICLE_KEY     = "NewsArticleID"
	ALL_ARTICLES_FEED    = "https://www.htafc.com/api/incrowd/getnewlistinformation?count=50"
	ONE_ARTICLE_FEED     = "https://www.htafc.com/api/incrowd/getnewsarticleinformation?id="
	WORKERS              = 3
	CRON_JOB_INTERVAL_MS = 300000
)

type Reader struct {
	dbCollection *mongo.Collection
	logger       *log.Logger
	httpClient   *http.Client
}

func NewReader(dbCollection *mongo.Collection, logger *log.Logger) *Reader {
	// timeout to prevent waiting too long for requests
	httpClient := &http.Client{
		Timeout: 4 * time.Second,
	}

	return &Reader{
		dbCollection: dbCollection,
		logger:       logger,
		httpClient:   httpClient,
	}
}

var wg sync.WaitGroup

func (r *Reader) RunCronFeedReader() error {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(CRON_JOB_INTERVAL_MS).Milliseconds().Do(r.feedNewsIntoDb)
	if err != nil {
		return err
	}
	s.StartAsync()
	return nil
}

func (r *Reader) feedNewsIntoDb() {
	newsListIds, err := r.getNewsList()
	if err != nil {
		r.logger.Printf("Error getting news list: %v", err)
		return
	}

	articleIDChan := make(chan int)

	for i := 0; i < WORKERS; i++ {
		go r.processArticles(articleIDChan, r.dbCollection)
	}

	for _, articleId := range newsListIds {
		wg.Add(1)
		articleIDChan <- articleId.NewsArticleID
	}

	close(articleIDChan)
	wg.Wait()
}

func (r *Reader) processArticles(articleIDChan <-chan int, dbCollection *mongo.Collection) {
	for articleID := range articleIDChan {
		article, err := r.getFullArticle(articleID)
		if err != nil {
			r.logger.Printf("Error getting article with id:%d and error: %v\n", articleID, err)
			wg.Done()
			continue
		}

		opts := options.Replace().SetUpsert(true)
		filter := bson.D{{Key: NEWS_ARTICLE_KEY, Value: articleID}}
		_, err = dbCollection.ReplaceOne(context.TODO(), filter, ConvertToMongoDB(article), opts)
		if err != nil {
			r.logger.Printf("Error saving article with id:%d and error: %v\n", articleID, err)
		}
		wg.Done()
	}
}

func (r *Reader) getNewsList() ([]NewsletterNewsItem, error) {
	response, err := http.Get(ALL_ARTICLES_FEED)
	if err != nil {
		r.logger.Printf("Error fetching the URL: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		r.logger.Printf("Error reading response: %v", err)
		return nil, err
	}

	var newsList NewListInformation
	err = xml.Unmarshal(body, &newsList)
	if err != nil {
		r.logger.Printf("Error unmarshaling XML: %v", err)
		return nil, err
	}
	return newsList.NewsletterNewsItems, nil
}

func (r *Reader) getFullArticle(articleID int) (*NewsArticleInformationXML, error) {
	url := fmt.Sprintf("%s%d", ONE_ARTICLE_FEED, articleID)

	response, err := http.Get(url)
	if err != nil {
		r.logger.Printf("Error fetching full article with id:%d and error: %v", articleID, err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		r.logger.Printf("Error reading response: %v", err)
		return nil, err
	}

	var article NewsArticleInformationXML
	err = xml.Unmarshal(body, &article)
	if err != nil {
		r.logger.Printf("Error unmarshaling XML: %v", err)
		return nil, err
	}
	return &article, nil
}
