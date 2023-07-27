package reader

import (
	"alibazlamit/feed-reader/database"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	response *http.Response
	err      error
}

func (c *MockHTTPClient) Get(url string) (*http.Response, error) {
	return c.response, c.err
}

func TestProcessArticles(t *testing.T) {
	articleIDChan := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(1)

	mockRepo := database.NewMockArticleRepository()
	mockLogger := log.New(nil, "", 0)
	mockHTTPClient := &MockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(strings.NewReader(`<NewsArticleInformation>
			<ClubName>TEST CITY</ClubName>
			<ClubWebsiteURL>test.com</ClubWebsiteURL>
			<NewsArticle>
			<ArticleURL>TEST URL</ArticleURL>
			<NewsArticleID>1</NewsArticleID>
			<PublishDate>2023-07-26 09:45:00</PublishDate>
			<Taxonomies>TEST News</Taxonomies>
			<TeaserText>TEST</TeaserText>
			<Subtitle/>
			<ThumbnailImageURL>test.png</ThumbnailImageURL>
			<Title>TEST</Title>
			<BodyText>test</BodyText>
			<GalleryImageURLs></GalleryImageURLs>
			<VideoURL/>
			<OptaMatchId/>
			<LastUpdateDate>2023-07-27 02:00:28</LastUpdateDate>
			<IsPublished>True</IsPublished>
			</NewsArticle>
			</NewsArticleInformation>`)),
		},
		err: nil,
	}

	reader := NewReader(mockRepo, mockLogger, mockHTTPClient)
	go reader.processArticles(articleIDChan, mockRepo, &wg)

	articleID := 123
	articleIDChan <- articleID
	close(articleIDChan)
	wg.Wait()

	assert.Equal(t, "TEST CITY", mockRepo.Articles[0].ClubName)
	assert.Equal(t, 1, mockRepo.Articles[0].NewsArticleID)

}

func TestGetFullArticle(t *testing.T) {
	mockRepo := database.NewMockArticleRepository()
	mockLogger := log.New(nil, "", 0)
	mockHTTPClient := &MockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(strings.NewReader(`<NewsArticleInformation>
			<ClubName>TEST CITY</ClubName>
			<ClubWebsiteURL>test.com</ClubWebsiteURL>
			<NewsArticle>
			<ArticleURL>TEST URL</ArticleURL>
			<NewsArticleID>1</NewsArticleID>
			<PublishDate>2023-07-26 09:45:00</PublishDate>
			<Taxonomies>TEST News</Taxonomies>
			<TeaserText>TEST</TeaserText>
			<Subtitle/>
			<ThumbnailImageURL>test.png</ThumbnailImageURL>
			<Title>TEST</Title>
			<BodyText>test</BodyText>
			<GalleryImageURLs></GalleryImageURLs>
			<VideoURL/>
			<OptaMatchId/>
			<LastUpdateDate>2023-07-27 02:00:28</LastUpdateDate>
			<IsPublished>True</IsPublished>
			</NewsArticle>
			</NewsArticleInformation>
			
			`)),
		},
		err: nil,
	}

	reader := NewReader(mockRepo, mockLogger, mockHTTPClient)

	articleID := 123
	artcl, err := reader.getFullArticle(articleID)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	assert.Equal(t, "TEST CITY", artcl.ClubName)
	assert.Equal(t, 1, artcl.NewsArticle.NewsArticleID)
}

func TestGetNewsList(t *testing.T) {
	mockRepo := database.NewMockArticleRepository()
	mockLogger := log.New(nil, "", 0)
	mockHTTPClient := &MockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(strings.NewReader(`<NewListInformation>
			<ClubName>TEST CITY</ClubName>
			<ClubWebsiteURL>test.com</ClubWebsiteURL>
			<NewsletterNewsItems>
			<NewsletterNewsItem>
			<ArticleURL>TEST URL</ArticleURL>
			<NewsArticleID>1</NewsArticleID>
			<PublishDate>2023-07-26 09:45:00</PublishDate>
			<Taxonomies>TEST News</Taxonomies>
			<TeaserText>TEST</TeaserText>
			<Subtitle/>
			<ThumbnailImageURL>test.png</ThumbnailImageURL>
			<Title>TEST</Title>
			<BodyText>test</BodyText>
			<GalleryImageURLs></GalleryImageURLs>
			<VideoURL/>
			<OptaMatchId/>
			<LastUpdateDate>2023-07-27 02:00:28</LastUpdateDate>
			<IsPublished>True</IsPublished>
			</NewsletterNewsItem>
			</NewsletterNewsItems>
			</NewListInformation>`)),
		},
		err: nil,
	}

	reader := NewReader(mockRepo, mockLogger, mockHTTPClient)

	newsList, err := reader.getNewsList()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(newsList) == 0 {
		t.Errorf("Expected non-empty news list, got empty")
	}
}

func TestRunCronFeedReader(t *testing.T) {
	mockRepo := database.NewMockArticleRepository()
	mockLogger := log.New(nil, "", 0)

	roundTripper := &customRoundTripper{
		responses: map[string]*http.Response{
			ALL_ARTICLES_FEED: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`<NewListInformation>
				<ClubName>TEST CITY</ClubName>
				<ClubWebsiteURL>test.com</ClubWebsiteURL>
				<NewsletterNewsItems>
				<NewsletterNewsItem>
				<ArticleURL>TEST URL</ArticleURL>
				<NewsArticleID>2</NewsArticleID>
				<PublishDate>2023-07-26 09:45:00</PublishDate>
				<Taxonomies>TEST News</Taxonomies>
				<TeaserText>TEST</TeaserText>
				<Subtitle/>
				<ThumbnailImageURL>test.png</ThumbnailImageURL>
				<Title>TEST</Title>
				<BodyText>test</BodyText>
				<GalleryImageURLs></GalleryImageURLs>
				<VideoURL/>
				<OptaMatchId/>
				<LastUpdateDate>2023-07-27 02:00:28</LastUpdateDate>
				<IsPublished>True</IsPublished>
				</NewsletterNewsItem>
				</NewsletterNewsItems>
				</NewListInformation>`)),
			},
			fmt.Sprintf("%s%d", ONE_ARTICLE_FEED, 2): &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`<NewsArticleInformation>
				<ClubName>TEST CITY</ClubName>
				<ClubWebsiteURL>test.com</ClubWebsiteURL>
				<NewsArticle>
				<ArticleURL>TEST URL</ArticleURL>
				<NewsArticleID>2</NewsArticleID>
				<PublishDate>2023-07-26 09:45:00</PublishDate>
				<Taxonomies>TEST News</Taxonomies>
				<TeaserText>TEST</TeaserText>
				<Subtitle/>
				<ThumbnailImageURL>test.png</ThumbnailImageURL>
				<Title>TEST</Title>
				<BodyText>test</BodyText>
				<GalleryImageURLs></GalleryImageURLs>
				<VideoURL/>
				<OptaMatchId/>
				<LastUpdateDate>2023-07-27 02:00:28</LastUpdateDate>
				<IsPublished>True</IsPublished>
				</NewsArticle>
				</NewsArticleInformation>
				
				`)),
			},
		},
	}

	mockHTTPClient := &http.Client{
		Transport: roundTripper,
	}

	reader := NewReader(mockRepo, mockLogger, mockHTTPClient)

	// Make the first request
	err := reader.RunCronFeedReader()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Wait for the cron jobs to run
	time.Sleep(2 * time.Second)

	assert.Equal(t, 1, len(mockRepo.Articles))
}

type customRoundTripper struct {
	responses map[string]*http.Response
}

func (c *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Get the response for the request URL.
	response, ok := c.responses[req.URL.String()]
	if !ok {
		return nil, fmt.Errorf("no response found for URL: %s", req.URL.String())
	}

	return response, nil
}
