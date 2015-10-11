package articles

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "gopkg.in/inconshreveable/log15.v2"
)

func scrapeUrl(a *Article) (err error) {
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
