package articles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"

	"github.com/PuerkitoBio/goquery"
	"github.com/boltdb/bolt"
	"github.com/gorilla/feeds"
	"github.com/rogierlommers/greedy/internal/common"
)

// BucketName is the name of the buckets
const BucketName = "articles"

var (
	db         *bolt.DB
	open       bool
	httpClient *http.Client
	s          Stats
)

// Article holds information about saved URL
type Article struct {
	ID          int
	URL         string
	Title       string
	Description string
	Added       time.Time
}

// NewClient creates http client with timeout
func NewClient() {
	httpClient = &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}
}

// Open creates database and opens it
func Open() (err error) {
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(common.Databasefile, 0600, config)
	if err != nil {
		log.Panicf("error creating bolt database: %s", err)
	}

	// create initial bucket (if not exists)
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	log.Infof("bucket initialized with %d records", count())
	open = true
	return nil
}

//Close closes database
func Close() {
	open = false
	db.Close()
}

func getArticles(amount int) (articleList []Article) {
	articleList = make([]Article, 0, 0)

	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(BucketName)).Cursor()
		var x = 0
		for _, v := c.Last(); x != amount; _, v = c.Prev() {
			var a *Article
			a, err := decode(v)
			if err != nil {
				return err
			}
			articleList = append(articleList, *a)
			x++
		}
		return nil
	})
	return articleList
}

// DisplayRSS generates RSS feed and writes to ResponseWriter
func DisplayRSS(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "your greedy's personal rss feed",
		Link:        &feeds.Link{Href: common.FeedsLink},
		Description: "Saved pages, all in one RSS feed",
		Author:      &feeds.Author{common.FeedsAuthorName, common.FeedsAuthorEmail},
		Created:     now,
	}

	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(BucketName)).Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var a *Article
			a, err := decode(v)
			if err != nil {
				return err
			}

			newItem := feeds.Item{
				Title:       a.Title,
				Link:        &feeds.Link{Href: a.URL},
				Description: a.Description,
				Created:     a.Added,
				Id:          strconv.Itoa(a.ID),
			}
			feed.Add(&newItem)
		}

		// update stats
		s.setLastCrawler(r.Referer())
		s.incCrawlCount()
		return nil
	})

	rss, err := feed.ToAtom()
	if err != nil {
		log.Errorf("error while generating RSS feed: %s", err)
		return
	}
	spew.Dump(r)
	w.Write([]byte(rss))
}

func count() (amount int) {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(BucketName)).Cursor()
		amount = c.Bucket().Stats().KeyN
		return nil
	})
	return amount
}

func decode(data []byte) (*Article, error) {
	var a *Article
	err := json.Unmarshal(data, &a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Scrape gathers information about new article
func (a *Article) Scrape() error {
	// time function duration
	start := time.Now()
	log.Infof("start scraping article [id: %d] [url: %s]", a.ID, a.URL)

	// init goquery
	doc, err := goquery.NewDocument(a.URL)
	if err != nil {
		return err
	}

	// start scraping page title
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		pageTitle := s.Find("title").Text()
		a.Title = strings.TrimSpace(pageTitle)
	})

	// now get meta description field
	// doc.Find("meta").Each(func(i int, s *goquery.Selection) {
	// 	if name, _ := s.Attr("name"); strings.EqualFold(name, "description") {
	// 		description, _ := s.Attr("content")
	// 		a.Description = strings.TrimSpace(description)
	// 	}
	// })

	// HERE WE SHOULD DOWNLOAD ORIGINAL HTML
	var sourceHTML string
	resp, err := httpClient.Get(a.URL)
	if err != nil {
		sourceHTML = "error while fetching: " + err.Error()
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sourceHTML = "error while fetching: " + err.Error()
	} else {
		sourceHTML = string(respByte)
	}
	a.Description = sourceHTML

	// debugging info
	elapsed := time.Since(start)
	log.Infof("scraping done [id: %d] [title: %q] [elapsed: %s]", a.ID, a.Title, elapsed)
	return nil
}

func (a *Article) encode() ([]byte, error) {
	enc, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// Save saves new article into db
func (a *Article) Save() error {
	if !open {
		return fmt.Errorf("db must be opened before saving")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		articles := tx.Bucket([]byte(BucketName))

		// Generate ID for the article.
		id, _ := articles.NextSequence()
		log.Infof("new sequence article: %d", id)
		a.ID = int(id)

		// scrape
		err := a.Scrape()
		if err != nil {
			log.Errorf("scraping error: %s", err)
		}

		enc, err := a.encode()
		if err != nil {
			return fmt.Errorf("could not encode article: %s", err)
		}

		err = articles.Put(itob(a.ID), enc)
		return err
	})
	return err
}
