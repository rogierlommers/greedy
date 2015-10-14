package articles

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"

	"github.com/PuerkitoBio/goquery"
	"github.com/boltdb/bolt"
	"github.com/rogierlommers/greedy/internal/common"
)

const BucketName = "articles"

var (
	db   *bolt.DB
	open bool
)

type Article struct {
	ID          uint64
	Url         string
	Added       time.Time
	Title       string
	Description string
}

func Open() (err error) {
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(common.Databasefile, 0600, config)
	if err != nil {
		log.Crit("error creating bolt database", "message", err)
	}
	open = true
	return nil
}

func Close() {
	open = false
	db.Close()
}

func (a *Article) GetId() (newId uint64, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		//articles, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		newId, err := tx.Bucket([]byte(BucketName)).NextSequence()
		if err != nil {
			return fmt.Errorf("error fetching new ID: %s", err)
		}
		log.Info("next sequence", "sequence", newId)
		return err
	})
	return newId, nil
}

func (a *Article) Save() error {
	if !open {
		return fmt.Errorf("db must be opened before saving")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		articles, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}
		enc, err := a.encode()
		if err != nil {
			return fmt.Errorf("could not encode article %s:", err)
		}

		// finally save article
		err = a.Scrape()
		if err != nil {
			log.Error("scraping error")
		}

		err = articles.Put(uint64ToByte(a.ID), enc)
		return err
	})
	return err
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
	log.Debug("scrape information", "title", a.Title, "description", a.Description)

	elapsed := time.Since(start)
	log.Info("scrape information", "time elapsed", elapsed.String())
	return nil

}

func getArticles() (articles []Article) {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(BucketName)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var a *Article
			a, err := decode(v)
			if err != nil {
				return err
			}

			log.Debug("fetching article", "id", "bla")
			article := Article{}
			article.ID = a.ID
			article.Title = a.Title
			article.Url = a.Url
			article.Added = a.Added
			article.Description = a.Description
			articles = append(articles, article)

		}
		return nil
	})
	return articles
}

func (a *Article) encode() ([]byte, error) {
	enc, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func decode(data []byte) (*Article, error) {
	var a *Article
	err := json.Unmarshal(data, &a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// ---------------------------------------------------------------------------------------
// still unused
// ---------------------------------------------------------------------------------------

//func GetArticle(id string) (*Article, error) {
//	if !open {
//		return nil, fmt.Errorf("db must be opened before saving!")
//	}
//	var a *Article
//	err := db.View(func(tx *bolt.Tx) error {
//		var err error
//		b := tx.Bucket([]byte(BucketName))
//		k := []byte(id)
//		a, err = decode(b.Get(k))
//		if err != nil {
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		fmt.Printf("Could not get Person ID %s", id)
//		return nil, err
//	}
//	return a, nil
//}
