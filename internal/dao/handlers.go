package dao

import (
	"encoding/base64"
	"html/template"
	"net/http"
	"path"
	"sort"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

func StatsHandler(database *ReadingListRecords) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		glog.Info("stats page")
		spew.Dump(database)
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
		vars := mux.Vars(r)
		base64url := vars["base64url"]

		urlByteArray, err := base64.StdEncoding.DecodeString(base64url)
		if err != nil {
			glog.Errorf("error decoding url -> %s", err)
			return
		}

		url := string(urlByteArray[:])

		AddRecord(database, url)
		glog.Infof("add url #%d --> [%s]: ", len(database.Records), url[0:60])
		w.Write([]byte("url added..."))
	}
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	glog.Info("index page")
	glog.Info(r)
	fp := path.Join("static", "templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO extract serverlocation from header
	obj := map[string]string{"serverLocation": "http://localhost:8080"}

	if err := tmpl.Execute(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
