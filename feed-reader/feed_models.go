package reader

import (
	"encoding/xml"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status string

const (
	Success Status = "success"
	Failure Status = "failure"
)

type NewListInformation struct {
	ClubName            string               `xml:"ClubName"`
	ClubWebsiteURL      string               `xml:"ClubWebsiteURL"`
	NewsletterNewsItems []NewsletterNewsItem `xml:"NewsletterNewsItems>NewsletterNewsItem"`
}

type NewsletterNewsItem struct {
	NewsArticleID int  `xml:"NewsArticleID"`
	IsPublished   bool `xml:"IsPublished"`
}

type NewsArticleInformationXML struct {
	ClubName       string      `xml:"ClubName"`
	ClubWebsiteURL string      `xml:"ClubWebsiteURL"`
	NewsArticle    NewsArticle `xml:"NewsArticle"`
}

type NewsArticle struct {
	ArticleURL        string     `xml:"ArticleURL"`
	NewsArticleID     int        `xml:"NewsArticleID"`
	PublishDate       CustomTime `xml:"PublishDate"`
	Taxonomies        string     `xml:"Taxonomies"`
	TeaserText        string     `xml:"TeaserText"`
	Subtitle          string     `xml:"Subtitle"`
	ThumbnailImageURL string     `xml:"ThumbnailImageURL"`
	Title             string     `xml:"Title"`
	BodyText          string     `xml:"BodyText"`
	GalleryImageURLs  string     `xml:"GalleryImageURLs"`
	VideoURL          string     `xml:"VideoURL"`
	OptaMatchID       string     `xml:"OptaMatchId"`
	LastUpdateDate    CustomTime `xml:"LastUpdateDate"`
	IsPublished       bool       `xml:"IsPublished"`
}

// Flattened structure for MongoDB
type NewsArticleInformationMongoDB struct {
	ClubName          string             `bson:"clubName" json:"clubName"`
	ClubWebsiteURL    string             `bson:"clubWebsiteURL" json:"-"`
	ArticleURL        string             `bson:"url" json:"url"`
	NewsArticleID     int                `bson:"NewsArticleID" json:"-"`
	PublishDate       time.Time          `bson:"publishDate" json:"published"`
	Taxonomies        string             `bson:"taxonomies" json:"-"`
	TeaserText        string             `bson:"teaser" json:"teaser"`
	Subtitle          string             `bson:"subtitle" json:"-"`
	ThumbnailImageURL string             `bson:"imageUrl" json:"imageUrl"`
	Title             string             `bson:"title" json:"title"`
	BodyText          string             `bson:"content" json:"content"`
	GalleryImageURLs  string             `bson:"galleryUrls" json:"galleryUrls"`
	VideoURL          string             `bson:"videoUrl" json:"videoUrl"`
	OptaMatchID       string             `bson:"optaMatchId" json:"optaMatchId"`
	LastUpdateDate    time.Time          `bson:"lastUpdateDate" json:"-"`
	IsPublished       bool               `bson:"published" json:"-"`
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

type NewsArticleResponse struct {
	Status   string                        `json:"status"`
	Data     NewsArticleInformationMongoDB `json:"data"`
	Metadata struct {
		CreatedAt string `json:"createdAt"`
	} `json:"metadata"`
	Error string `json:"error,omitempty"`
}

type NewsArticlesResponse struct {
	Status   string                          `json:"status"`
	Data     []NewsArticleInformationMongoDB `json:"data"`
	Metadata struct {
		CreatedAt  string `json:"createdAt"`
		TotalItems int    `json:"totalItems"`
		Sort       string `json:"sort"`
	} `json:"metadata"`
	Error string `json:"error,omitempty"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const layout = "2006-01-02 15:04:05"
	var v string
	d.DecodeElement(&v, &start)
	t, err := time.Parse(layout, v)
	if err != nil {
		return err
	}
	*ct = CustomTime{t}
	return nil
}

func ConvertToMongoDB(newsArticleInfo *NewsArticleInformationXML) *NewsArticleInformationMongoDB {
	newsArticleInfoMongoDB := NewsArticleInformationMongoDB{
		ClubName:          newsArticleInfo.ClubName,
		ClubWebsiteURL:    newsArticleInfo.ClubWebsiteURL,
		ArticleURL:        newsArticleInfo.NewsArticle.ArticleURL,
		NewsArticleID:     newsArticleInfo.NewsArticle.NewsArticleID,
		PublishDate:       newsArticleInfo.NewsArticle.PublishDate.Time,
		Taxonomies:        newsArticleInfo.NewsArticle.Taxonomies,
		TeaserText:        newsArticleInfo.NewsArticle.TeaserText,
		Subtitle:          newsArticleInfo.NewsArticle.Subtitle,
		ThumbnailImageURL: newsArticleInfo.NewsArticle.ThumbnailImageURL,
		Title:             newsArticleInfo.NewsArticle.Title,
		BodyText:          newsArticleInfo.NewsArticle.BodyText,
		GalleryImageURLs:  newsArticleInfo.NewsArticle.GalleryImageURLs,
		VideoURL:          newsArticleInfo.NewsArticle.VideoURL,
		OptaMatchID:       newsArticleInfo.NewsArticle.OptaMatchID,
		LastUpdateDate:    newsArticleInfo.NewsArticle.LastUpdateDate.Time,
		IsPublished:       newsArticleInfo.NewsArticle.IsPublished,
	}
	return &newsArticleInfoMongoDB
}
