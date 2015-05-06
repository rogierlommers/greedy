package dao

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"sort"
	"time"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/golang/glog"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

// base64 for about:blank pages, which we don't want to store
const AboutBlank = "YWJvdXQ6Ymxhbms="

func StatsHandler(database *ReadingListRecords) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fp := path.Join("static", "templates", "stats.html")
		tmpl, parseErr := template.ParseFiles(fp)
		if parseErr != nil {
			http.Error(w, parseErr.Error(), http.StatusInternalServerError)
			return
		}

		var jsonBytes []byte
		var jsonErr error

		jsonBytes, jsonErr = json.MarshalIndent(stats_api.GetStats(), "", "  ")

		var stats string
		if jsonErr != nil {
			stats = jsonErr.Error()
		} else {
			stats = string(jsonBytes)
		}

		obj := map[string]string{"statsmessage": stats}

		if templErr := tmpl.Execute(w, obj); templErr != nil {
			http.Error(w, templErr.Error(), http.StatusInternalServerError)
		}

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

		if isAboutBlank(base64url) {
			IndexPage(w, r)
		}

		urlByteArray, decodeErr := base64.StdEncoding.DecodeString(base64url)
		if decodeErr != nil {
			glog.Errorf("error decoding url -> %s", decodeErr)

			fp := path.Join("static", "templates", "error.html")
			tmpl, parseErr := template.ParseFiles(fp)
			if parseErr != nil {
				http.Error(w, parseErr.Error(), http.StatusInternalServerError)
				return
			}

			obj := map[string]string{"errormessage": decodeErr.Error()}

			if templErr := tmpl.Execute(w, obj); templErr != nil {
				http.Error(w, templErr.Error(), http.StatusInternalServerError)
			}
			return
		}

		addedUrl := string(urlByteArray[:])
		AddRecord(database, addedUrl)

		var logUrl = ""
		if len(addedUrl) < 60 {
			logUrl = addedUrl
		} else {
			logUrl = addedUrl[0:60]
		}
		glog.Infof("add url #%d --> [%s]", len(database.Records), logUrl)

		fp := path.Join("static", "templates", "confirmation.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		u, _ := url.Parse(addedUrl)
		obj := map[string]string{"url": u.Host}

		if err := tmpl.Execute(w, obj); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
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

func isAboutBlank(url string) bool {
	return url == AboutBlank
}

func getStringFromBase64(url string) (string, error) {
	// validates url and returns string of url
	glog.Info("check validiry")
	return "", nil
}
