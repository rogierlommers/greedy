package articles

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"

	"github.com/PuerkitoBio/goquery"
	"github.com/boltdb/bolt"
	"github.com/gorilla/feeds"
	"github.com/rogierlommers/greedy/internal/common"
)

const BucketName = "articles"

var (
	db   *bolt.DB
	open bool
)

type Article struct {
	ID          int
	Url         string
	Title       string
	Description string
	Added       time.Time
}

func Open() (err error) {
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(common.Databasefile, 0600, config)
	if err != nil {
		log.Crit("error creating bolt database", "message", err)
	}

	// create initial bucket (if not exists)
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	log.Info("bucket", "amount", count())
	open = true
	return nil
}

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
				Link:        &feeds.Link{Href: a.Url},
				Description: a.Description,
				Created:     a.Added,
				Id:          strconv.Itoa(a.ID),
			}
			feed.Add(&newItem)
		}
		return nil
	})

	rss, err := feed.ToAtom()
	if err != nil {
		log.Error("error generation RSS feed", "message", err)
		return
	}
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

func (a *Article) Scrape() error {
	// time function duration
	start := time.Now()
	log.Info("start scraping article", "time", start, "id", a.ID, "url", a.Url)

	// init goquery
	doc, err := goquery.NewDocument(a.Url)
	if err != nil {
		return err
	}

	// start scraping page title
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		pageTitle := s.Find("title").Text()
		a.Title = strings.TrimSpace(pageTitle)
	})

	// now get meta description field
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); strings.EqualFold(name, "description") {
			description, _ := s.Attr("content")
			a.Description = strings.TrimSpace(description)
		}
	})

	// debugging info
	elapsed := time.Since(start)
	log.Debug("scraped", "id", a.ID, "title", a.Title, "description", a.Description, "elapsed", elapsed)
	return nil
}

func (a *Article) encode() ([]byte, error) {
	enc, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (a *Article) Save() error {
	if !open {
		return fmt.Errorf("db must be opened before saving")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		articles := tx.Bucket([]byte(BucketName))

		// Generate ID for the article.
		id, _ := articles.NextSequence()
		log.Info("article", "sequence", id)
		a.ID = int(id)

		// scrape
		err := a.Scrape()
		if err != nil {
			log.Error("scraping error")
		}

		enc, err := a.encode()
		if err != nil {
			return fmt.Errorf("could not encode article %s:", err)
		}

		err = articles.Put(itob(a.ID), enc)
		return err
	})
	return err
}
