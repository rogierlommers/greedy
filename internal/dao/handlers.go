package dao

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/golang/glog"
	"github.com/gorilla/feeds"
	"github.com/rogierlommers/go-read/internal/render"
)

func StatsHandler(database *ReadingListRecords) http.HandlerFunc {
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

func GenerateRSS(database *ReadingListRecords) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sort.Sort(sort.Reverse(ById(database.Records)))

		now := time.Now()
		feed := &feeds.Feed{
			Title:       "Go-read",
			Link:        &feeds.Link{Href: "http://bla.com"},
			Description: "personal RSS feed with articles to be read",
			Author:      &feeds.Author{"Rogier Lommers", "rogier@lommers.org"},
			Created:     now,
		}

		for _, value := range database.Records {
			newItem := feeds.Item{Title: value.URL,
				Link: &feeds.Link{Href: value.URL},
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

func AddArticle(database *ReadingListRecords) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParam := r.FormValue("url")
		glog.Info("len(queryParam): ", len(queryParam))
		if len(queryParam) == 0 || queryParam == "about:blank" {
			renderObject := map[string]string{"errorMessage": "unable to insert empty or about:blank page"}
			render.DisplayPage(w, r, renderObject, "error.html")
			return
		}

		amount := AddRecord(database, queryParam)
		addedUrl := logAddedUrl(queryParam, database)

		renderObject := map[string]string{"url": addedUrl, "amount": strconv.Itoa(amount)}
		render.DisplayPage(w, r, renderObject, "confirmation.html")

	}
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	// TODO extract serverlocation from header
	renderObject := map[string]string{"serverLocation": "http://localhost:8080"}
	render.DisplayPage(w, r, renderObject, "index.html")
}

func logAddedUrl(addedUrl string, database *ReadingListRecords) (rogier string) {
	var logUrl = ""
	if len(addedUrl) < 60 {
		logUrl = addedUrl
	} else {
		logUrl = addedUrl[0:60]
	}
	glog.Infof("add url #%d --> [%s]", len(database.Records), logUrl)
	return logUrl
}
