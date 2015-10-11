package articles

import (
	"database/sql"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "gopkg.in/inconshreveable/log15.v2"
)

func scrapeArticle(id int64) {
	// time function duration
	start := time.Now()

	// storedArticle contains information stored in db which need to be updated through scraping
	storedArticle := getArticleById(id)

	// init goquery
	doc, err := goquery.NewDocument(storedArticle.Url.String)
	log.Info("start scraping article", "url", storedArticle.Url.String)

	if err != nil {
		log.Error("error while scraping article", "article id", id, "message", err)
		return
	}

	// start scraping page title
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		pageTitle := s.Find("title").Text()
		storedArticle.Name = sql.NullString{String: strings.TrimSpace(pageTitle), Valid: true}
	})

	// now get meta description field
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); strings.EqualFold(name, "description") {
			description, _ := s.Attr("content")
			storedArticle.Description = sql.NullString{String: strings.TrimSpace(description), Valid: true}
		}
	})

	// if unable to scrape title, then use url
	if len(storedArticle.Name.String) == 0 {
		storedArticle.Name.String = storedArticle.Url.String
	}

	// if unable to scrape description, then use default text
	if len(storedArticle.Description.String) == 0 {
		storedArticle.Description.String = noDescription
	}

	// debugging info
	log.Debug("scrape information", "title", storedArticle.Name.String, "title length", len(storedArticle.Name.String))
	log.Debug("scrape information", "description", storedArticle.Description.String, "description length", len(storedArticle.Description.String))

	// after succesfull scraping, add page title (and more?) to article in db
	updateArticle(db, storedArticle)

	elapsed := time.Since(start)
	log.Info("scrape information", "stored article id", storedArticle.ID.Int64, "time elapsed", elapsed.String())
}
