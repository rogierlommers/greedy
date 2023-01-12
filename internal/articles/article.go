package articles

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/badoux/goscraper"
	"github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
	"github.com/gorilla/feeds"
	"github.com/rogierlommers/greedy/internal/common"
)

const (
	bucketName      = "articles"
	keep            = 100   // amount of records to keep on disk
	numberInRSS     = 25    // amount of records to display in feed
	scheduleCleanup = 86400 // 1 day
)

var (
	db   *bolt.DB
	open bool
	s    Stats
)

// Article holds information about saved URL
type Article struct {
	ID          int
	URL         string
	Title       string
	Description string
	Added       time.Time
}

// Open creates database and opens it
func Open() (err error) {
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(common.Databasefile, 0600, config)
	if err != nil {
		logrus.Panicf("error creating bolt database: %s", err)
	}

	// create initial bucket (if not exists)
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	logrus.Infof("bucket initialized with %d records", count())
	open = true
	return nil
}

// Close closes database
func Close() {
	open = false
	db.Close()
}

func getArticles(amount int) (articleList []Article) {
	articleList = make([]Article, 0)

	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucketName)).Cursor()
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
		Author: &feeds.Author{
			Name:  common.FeedsAuthorName,
			Email: common.FeedsAuthorEmail,
		},
		Created: now,
	}

	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucketName)).Cursor()
		count := 0
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			if count >= numberInRSS {
				break
			}

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
			count++
		}

		// update stats
		s.setLastCrawler(r.UserAgent())
		s.incCrawlCount()
		return nil
	})

	rss, err := feed.ToAtom()
	if err != nil {
		logrus.Errorf("error while generating RSS feed: %s", err)
		return
	}
	w.Write([]byte(rss))
}

func count() (amount int) {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucketName)).Cursor()
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
	logrus.Infof("start scraping article [id: %d] [url: %s]", a.ID, a.URL)

	s, err := goscraper.Scrape(a.URL, 5)
	if err != nil {
		a.Title = fmt.Sprintf("[Greedy] scrape failed: %q", a.URL)
		a.Description = fmt.Sprintf("Scraping failed for url %q", a.URL)
		logrus.Errorf("scrape error: %s", err)
	} else {
		a.Title = fmt.Sprintf("[Greedy] %s", s.Preview.Title)
		a.Description = s.Preview.Description
	}

	// debugging info
	elapsed := time.Since(start)
	logrus.Infof("scraping done [id: %d] [title: %q] [elapsed: %s]", a.ID, a.Title, elapsed)
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
		articles := tx.Bucket([]byte(bucketName))

		// Generate ID for the article.
		id, _ := articles.NextSequence()
		logrus.Infof("new sequence article: %d", id)
		a.ID = int(id)

		// scrape
		err := a.Scrape()
		if err != nil {
			logrus.Errorf("scraping error: %s", err)
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

func cleanUp(numberToKeep int) int {
	var (
		count   = 0
		deleted = 0
	)

	db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucketName)).Cursor()

		for k, _ := c.Last(); k != nil; k, _ = c.Prev() {
			count++
			if count > numberToKeep {
				err := c.Delete()
				if err != nil {
					logrus.Errorf("error deleting record while cleanup: %q", err)
				} else {
					deleted++
				}
			}
		}
		return nil
	})
	return deleted
}

// ScheduleCleanup removes old records
func ScheduleCleanup() {
	go func() {
		logrus.Infof("scheduled cleanup, every %d seconds, remove more than %d records", scheduleCleanup, keep)
		for {
			deleted := cleanUp(keep)
			logrus.Infof("deleted %d records from database", deleted)
			time.Sleep(scheduleCleanup * time.Second)
		}
	}()
}
