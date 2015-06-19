package handlers

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"net/http"
	"net/url"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/glog"
	"github.com/gorilla/feeds"
	"github.com/rogierlommers/go-read/internal/common"
	"github.com/rogierlommers/go-read/internal/dao"
	"github.com/rogierlommers/go-read/internal/render"
)

func StatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var stats = "<table><tr><th>Title</th><th>Added</th></tr>"
		articles := dao.GetLastArticles(db)

		for _, value := range articles {
			stats += "<tr><td>" + getHostnameFromUrl(value.Url.String) + "</td><td>" + humanize.Time(value.Created) + "</td></tr>"
		}

		renderObject := map[string]string{"message": stats}
		render.DisplayPage(w, r, renderObject, "stats.html")

	}

}

func ExportCSV(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		articles := dao.GetLastArticles(db)

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
			Title:       "Go-read",
			Link:        &feeds.Link{Href: "http://read.lommers.org/"},
			Description: "personal feed with saved articles",
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
				Id:          getMD5Hash(value.Url.String + value.Created.String()),
			}
			feed.Add(&newItem)
		}

		rss, err := feed.ToAtom()
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
		glog.Infof("add hostname %s, id: %d", getHostnameFromUrl(queryParam), insertedId)

		// start routine which scrapes url
		go dao.ScrapeArticle(db, insertedId)

		// finally output confirmation page
		renderObject := map[string]string{"hostname": getHostnameFromUrl(queryParam)}
		render.DisplayPage(w, r, renderObject, "confirmation.html")
	}
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	renderObject := map[string]string{"serverLocation": "http://read.lommers.org", "buildversion": common.BuildDate}
	render.DisplayPage(w, r, renderObject, "index.html")
}

func getHostnameFromUrl(addedUrl string) (hostname string) {
	u, err := url.Parse(addedUrl)
	if err != nil {
		glog.Error(err)
	}
	return u.Host
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
