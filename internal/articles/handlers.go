package articles

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/render"
)

// AddArticle stores new article into database
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
		URL:   queryParam,
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

// IndexPage writes IndexPage to ResponseWriter
func IndexPage(w http.ResponseWriter, r *http.Request) {
	host := fmt.Sprint("http://", r.Host) // needed to build bookmarklet
	renderObject := map[string]interface{}{
		"IsLandingPage":  "true",
		"serverLocation": host,
		"buildversion":   common.BuildDate,
		"statsHTML":      template.HTML(getStatsHTML()),
		"amount":         count(),
	}
	render.DisplayPage(w, r, renderObject)
}

// ExportCSV sends csv file to ResponseWriter
func ExportCSV(w http.ResponseWriter, r *http.Request) {
	articles := getArticles(count())

	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)

	for _, value := range articles {
		record := []string{value.URL}
		wr.Write(record)
	}

	wr.Flush()
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=go-read-articles.csv")
	w.Write(b.Bytes())
}

func getStatsHTML() string {
	var stats = "<table><tr><th>Title</th><th>Added</th></tr>"

	articles := getArticles(100)
	for _, value := range articles {
		stats += "<tr><td>" + getHostnameFromUrl(value.URL) + "</td><td>" + humanize.Time(value.Added) + "</td></tr>"
	}
	return stats
}
