package handlers

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gorilla/feeds"
	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/dao"
	"github.com/rogierlommers/greedy/internal/render"
	log "gopkg.in/inconshreveable/log15.v2"
)

func StatsHandler(db *sql.DB) http.HandlerFunc {
	// http://julianyap.com/2013/09/23/using-anonymous-structs-to-pass-data-to-templates-in-golang.html
	return func(w http.ResponseWriter, r *http.Request) {
		amount := dao.GetNumberOfRecords(db)

		var stats = "<table><tr><th>Title</th><th>Added</th></tr>"
		articles := dao.GetLastArticles(db, 0)
		for _, value := range articles {
			stats += "<tr><td>" + getHostnameFromUrl(value.Url.String) + "</td><td>" + humanize.Time(value.Created) + "</td></tr>"
		}

		renderObject := map[string]interface{}{
			"IsStatsPage": "true",
			"amount":      strconv.Itoa(amount),
			"html":        template.HTML(stats),
		}
		render.DisplayPage(w, r, renderObject)
	}
}

func ExportCSV(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		articles := dao.GetLastArticles(db, 0)

		b := &bytes.Buffer{}
		wr := csv.NewWriter(b)

		for _, value := range articles {
			record := []string{value.Url.String}
			wr.Write(record)
		}

		wr.Flush()
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=go-read-articles.csv")
		w.Write(b.Bytes())
	}
}

func GenerateRSS(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		feed := &feeds.Feed{
			Title:       "your greedy's personal rss feed",
			Link:        &feeds.Link{Href: common.FeedsLink},
			Description: "discussion about tech, footie, photos",
			Author:      &feeds.Author{common.FeedsAuthorName, common.FeedsAuthorEmail},
			Created:     now,
		}

		var articles []dao.ArticleStruct
		articles = dao.GetLastArticles(db, 100)

		for _, value := range articles {
			newItem := feeds.Item{
				Title:       value.Name.String,
				Link:        &feeds.Link{Href: value.Url.String},
				Description: value.Description.String,
				Created:     value.Created,
				Id:          getMD5Hash(value.Url.String + value.Created.String()),
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
}

func AddArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParam := r.FormValue("url")
		if len(queryParam) == 0 || queryParam == "about:blank" {

			renderObject := map[string]interface{}{
				"IsErrorPage":  "true",
				"errorMessage": "unable to insert empty or about:blank page",
			}
			render.DisplayPage(w, r, renderObject)
			return
		}

		insertedId := dao.SaveArticle(db, queryParam)
		log.Info("adding article", "hostname", getHostnameFromUrl(queryParam), "article id", insertedId)

		// start routine which scrapes url
		go dao.ScrapeArticle(db, insertedId)

		// finally output confirmation page
		renderObject := map[string]interface{}{
			"IsConfirmation": "true",
			"hostname":       getHostnameFromUrl(queryParam),
		}
		render.DisplayPage(w, r, renderObject)
	}
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
