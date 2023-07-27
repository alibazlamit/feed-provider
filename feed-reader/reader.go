package reader

import (
	"alibazlamit/feed-reader/database"
	"alibazlamit/feed-reader/models"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

const (
	NEWS_ARTICLE_KEY     = "NewsArticleID"
	ALL_ARTICLES_FEED    = "https://www.htafc.com/api/incrowd/getnewlistinformation?count=50"
	ONE_ARTICLE_FEED     = "https://www.htafc.com/api/incrowd/getnewsarticleinformation?id="
	WORKERS              = 5
	CRON_JOB_INTERVAL_MS = 300000
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type Reader struct {
	db         database.ArticleRepository
	logger     *log.Logger
	httpClient HTTPClient
}

func NewReader(db database.ArticleRepository, logger *log.Logger, httpClient HTTPClient) *Reader {

	return &Reader{
		db:         db,
		logger:     logger,
		httpClient: httpClient,
	}
}

func (r *Reader) RunCronFeedReader() error {
	// run a cron every interval in milliseconds
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(CRON_JOB_INTERVAL_MS).Milliseconds().Do(r.feedNewsIntoDb)
	if err != nil {
		return err
	}
	s.StartAsync()
	return nil
}

func (r *Reader) feedNewsIntoDb() {
	var wg sync.WaitGroup
	//read news feed
	newsListIds, err := r.getNewsList()
	if err != nil {
		r.logger.Printf("Error getting news list: %v", err)
		return
	}

	//create a buffered channel of the number of workers set
	articleIDChan := make(chan int, WORKERS)

	// sync process all articles from feed at the same time
	for i := 0; i < WORKERS; i++ {
		go r.processArticles(articleIDChan, r.db, &wg)
	}

	for _, articleId := range newsListIds {
		wg.Add(1)
		articleIDChan <- articleId.NewsArticleID
	}

	close(articleIDChan)
	wg.Wait()
}

func (r *Reader) processArticles(articleIDChan <-chan int, db database.ArticleRepository, wg *sync.WaitGroup) {
	for articleID := range articleIDChan {
		article, err := r.getFullArticle(articleID)
		if err != nil {
			r.logger.Printf("Error getting article with id:%d and error: %v\n", articleID, err)
			wg.Done()
			continue
		}
		err = r.db.AddOrUpdateArticle(articleID, article)
		if err != nil {
			r.logger.Printf("Error saving article with id:%d and error: %v\n", articleID, err)
		}
		wg.Done()
	}
}

// reading from feed and transforming xml into structs
func (r *Reader) getNewsList() ([]models.NewsletterNewsItem, error) {
	response, err := r.httpClient.Get(ALL_ARTICLES_FEED)
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

	var newsList models.NewListInformation
	err = xml.Unmarshal(body, &newsList)
	if err != nil {
		r.logger.Printf("Error unmarshaling XML: %v", err)
		return nil, err
	}
	return newsList.NewsletterNewsItems, nil
}

// reading from feed and transforming xml into structs
func (r *Reader) getFullArticle(articleID int) (*models.NewsArticleInformationXML, error) {
	url := fmt.Sprintf("%s%d", ONE_ARTICLE_FEED, articleID)

	response, err := r.httpClient.Get(url)
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

	var article models.NewsArticleInformationXML
	err = xml.Unmarshal(body, &article)
	if err != nil {
		r.logger.Printf("Error unmarshaling XML: %v", err)
		return nil, err
	}
	return &article, nil
}
