package articles

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"time"

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
		ID:    getMD5Hash(queryParam),
		Url:   queryParam,
		Added: time.Now(),
	}

	err := newArticle.Save()
	if err != nil {
		log.Warn("error saving article", "hostname", getHostnameFromUrl(queryParam), "article id", newArticle.ID)
	}
	log.Info("article added", "hostname", getHostnameFromUrl(queryParam), "article id", newArticle.ID)

	// start routine which scrapes url
	go scrapeArticle(newArticle.ID)

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

func getHostnameFromUrl(addedUrl string) (hostname string) {
	u, err := url.Parse(addedUrl)
	if err != nil {
		log.Error("error looking up hostname from url", "url", addedUrl, "message", err)
	}
	return u.Host
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
