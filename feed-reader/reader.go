package reader

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	NEWS_ARTICLE_KEY  = "NewsArticleID"
	ALL_ARTICLES_FEED = "https://www.htafc.com/api/incrowd/getnewlistinformation?count=50"
	ONE_ARTICLE_FEED  = "https://www.htafc.com/api/incrowd/getnewsarticleinformation?id="
	WORKERS           = 10
)

var wg sync.WaitGroup

func RunCronFeedReader(dbCollection *mongo.Collection) error {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(5).Seconds().Do(feedNewsIntoDb, dbCollection)
	if err != nil {
		return err
	}
	s.StartAsync()
	return nil
}

func feedNewsIntoDb(dbCollection *mongo.Collection) {
	newsListIds, err := getNewsList()
	if err != nil {
		fmt.Println("Error getting news list:", err)
		return
	}

	articleIDChan := make(chan int)

	for i := 0; i < WORKERS; i++ {
		go processArticles(articleIDChan, dbCollection)
	}

	for _, articleId := range newsListIds {
		wg.Add(1)
		articleIDChan <- articleId.NewsArticleID
	}

	close(articleIDChan)
	wg.Wait()
}

func processArticles(articleIDChan <-chan int, dbCollection *mongo.Collection) {
	for articleID := range articleIDChan {
		article, err := getFullArticle(articleID)
		if err != nil {
			fmt.Printf("Error getting article with id:%d and error: %v\n", articleID, err)
			continue
		}

		opts := options.Replace().SetUpsert(true)
		filter := bson.D{{Key: NEWS_ARTICLE_KEY, Value: articleID}}
		_, err = dbCollection.ReplaceOne(context.TODO(), filter, ConvertToMongoDB(article), opts)
		if err != nil {
			fmt.Printf("Error saving article with id:%d and error: %v\n", articleID, err)
		}
		wg.Done()
	}
}

func getNewsList() ([]NewsletterNewsItem, error) {
	response, err := http.Get(ALL_ARTICLES_FEED)
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	var newsList NewListInformation
	err = xml.Unmarshal(body, &newsList)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return nil, err
	}
	return newsList.NewsletterNewsItems, nil
}

func getFullArticle(articleID int) (*NewsArticleInformationXML, error) {
	url := fmt.Sprintf("%s%d", ONE_ARTICLE_FEED, articleID)

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching full article with id:%d and error: %v", articleID, err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	var article NewsArticleInformationXML
	err = xml.Unmarshal(body, &article)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return nil, err
	}
	return &article, nil

}
