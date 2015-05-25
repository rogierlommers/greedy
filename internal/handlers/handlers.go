package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/golang/glog"
	"github.com/gorilla/feeds"
	"github.com/rogierlommers/go-read/internal/common"
	"github.com/rogierlommers/go-read/internal/dao"
	"github.com/rogierlommers/go-read/internal/render"
)

func StatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonBytes []byte
		var jsonErr error
		var stats string

		jsonBytes, jsonErr = json.MarshalIndent(stats_api.GetStats(), "", "  ")

		if jsonErr != nil {
			stats = jsonErr.Error()
		} else {
			stats = string(jsonBytes)
		}

		renderObject := map[string]string{"message": stats}
		render.DisplayPage(w, r, renderObject, "stats.html")
	}
}

func GenerateRSS(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//sort.Sort(sort.Reverse(ById(database.Records)))

		now := time.Now()
		feed := &feeds.Feed{
			Title:       "Go-read",
			Link:        &feeds.Link{Href: "http://bla.com"},
			Description: "personal RSS feed with articles to be read",
			Author:      &feeds.Author{"Rogier Lommers", "rogier@lommers.org"},
			Created:     now,
		}

		var articles []dao.ArticleStruct
		articles = dao.GetLastArticles(db)

		for _, value := range articles {
			newItem := feeds.Item{
				Title:       value.Name.String,
				Link:        &feeds.Link{Href: value.Url.String},
				Description: value.Description.String,
				Created:     value.Created,
			}
			feed.Add(&newItem)
		}

		rss, err := feed.ToRss()
		if err != nil {
			glog.Errorf("error creating RSS feed -> %s", err)
			return
		}
		w.Write([]byte(rss))
	}
}

func AddArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParam := r.FormValue("url")
		if len(queryParam) == 0 || queryParam == "about:blank" {
			renderObject := map[string]string{"errorMessage": "unable to insert empty or about:blank page"}
			render.DisplayPage(w, r, renderObject, "error.html")
			return
		}

		insertedId := dao.SaveArticle(db, queryParam)
		addedUrl := logAddedUrl(queryParam, insertedId)

		// start routine which scrapes url
		dao.ScrapeArticle(insertedId)

		// finally output confirmation page
		renderObject := map[string]string{"url": addedUrl, "amount": "1"}
		render.DisplayPage(w, r, renderObject, "confirmation.html")
	}
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	// TODO extract serverlocation from header
	renderObject := map[string]string{"serverLocation": "http://read.lommers.org", "buildversion": common.BuildDate}
	render.DisplayPage(w, r, renderObject, "index.html")
}

func logAddedUrl(addedUrl string, insertedId int64) (rogier string) {
	var logUrl = ""
	if len(addedUrl) < 60 {
		logUrl = addedUrl
	} else {
		logUrl = addedUrl[0:60]
	}
	glog.Infof("add url (id: %s) --> [%d]", logUrl, insertedId)
	return logUrl
}
