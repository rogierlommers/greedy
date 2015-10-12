package articles

import (
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/render"
	log "gopkg.in/inconshreveable/log15.v2"
)

func AddArticle(w http.ResponseWriter, r *http.Request) {
	queryParam := r.FormValue("url")
	if len(queryParam) == 0 || queryParam == "about:blank" {

		renderObject := map[string]interface{}{
			"IsErrorPage":  "true",
			"errorMessage": "unable to insert empty or about:blank page",
		}
		render.DisplayPage(w, r, renderObject)
		return
	}

	newArticle := Article{
		//ID:    getMD5Hash(queryParam),
		Url:   queryParam,
		Added: time.Now(),
	}

	err := newArticle.Save()
	if err != nil {
		log.Warn("error saving article", "hostname", getHostnameFromUrl(queryParam), "article id", newArticle.ID)
	}

	log.Info("article added", "hostname", getHostnameFromUrl(queryParam), "article id", newArticle.ID)

	// start routine which scrapes url
	err = newArticle.Scrape()
	if err != nil {
		log.Warn("error scraping article", "hostname", getHostnameFromUrl(queryParam), "article id", newArticle.ID)
	}

	// finally output confirmation page
	renderObject := map[string]interface{}{
		"IsConfirmation": "true",
		"hostname":       getHostnameFromUrl(queryParam),
	}
	render.DisplayPage(w, r, renderObject)
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	renderObject := map[string]interface{}{
		"IsLandingPage":  "true",
		"serverLocation": r.Host,
		"buildversion":   common.BuildDate,
	}
	render.DisplayPage(w, r, renderObject)
}

func GenerateRSS(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "your greedy's personal rss feed",
		Link:        &feeds.Link{Href: common.FeedsLink},
		Description: "discussion about tech, footie, photos",
		Author:      &feeds.Author{common.FeedsAuthorName, common.FeedsAuthorEmail},
		Created:     now,
	}

	articles := getArticles()

	for _, a := range articles {
		newItem := feeds.Item{
			Title:       a.Title,
			Link:        &feeds.Link{Href: a.Url},
			Description: a.Description,
			Created:     a.Added,
			Id:          a.ID,
		}
		feed.Add(&newItem)
	}

	rss, err := feed.ToAtom()
	if err != nil {
		log.Error("error generation RSS feed", "message", err)
		return
	}
	w.Write([]byte(rss))
}
