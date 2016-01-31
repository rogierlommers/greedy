package articles

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/render"
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
		Url:   queryParam,
		Added: time.Now(),
	}

	err := newArticle.Save()
	if err != nil {
		log.Warn("error saving article", "hostname", getHostnameFromUrl(queryParam), "id", newArticle.ID)
	}

	// finally output confirmation page
	renderObject := map[string]interface{}{
		"IsConfirmation": "true",
		"hostname":       getHostnameFromUrl(queryParam),
	}
	render.DisplayPage(w, r, renderObject)
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	host := fmt.Sprint("http://", r.Host) // needed to build bookmarklet
	renderObject := map[string]interface{}{
		"IsLandingPage":  "true",
		"serverLocation": host,
		"buildversion":   common.BuildDate,
	}
	render.DisplayPage(w, r, renderObject)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	var stats = "<table><tr><th>Title</th><th>Added</th></tr>"

	articles := getArticles(100)
	for _, value := range articles {
		stats += "<tr><td>" + getHostnameFromUrl(value.Url) + "</td><td>" + humanize.Time(value.Added) + "</td></tr>"
	}

	renderObject := map[string]interface{}{
		"IsStatsPage": "true",
		"amount":      strconv.Itoa(count()),
		"html":        template.HTML(stats),
	}
	render.DisplayPage(w, r, renderObject)
}

func ExportCSV(w http.ResponseWriter, r *http.Request) {
	articles := getArticles(count())

	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)

	for _, value := range articles {
		record := []string{value.Url}
		wr.Write(record)
	}

	wr.Flush()
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=go-read-articles.csv")
	w.Write(b.Bytes())
}
