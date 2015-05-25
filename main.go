package main

import (
	"database/sql"
	"flag"
	"net/http"
	"os"
	"strconv"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	"github.com/rogierlommers/go-read/internal/common"
	"github.com/rogierlommers/go-read/internal/dao"
	"github.com/rogierlommers/go-read/internal/handlers"
)

// TODO
// https://github.com/siadat/eton
// sqlite branche wijziging gemaakt op laptop
// injected by the build process
var builddate = "unknown build date"

// read flags
var databasefile = flag.String("databasefile", "articles.db", "sqllite file where items are stored")
var port = flag.Int("port", 8080, "http listener port")

func init() {
	flag.Parse()
	flag.Lookup("alsologtostderr").Value.Set("true")
}

func log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	defer glog.Flush()

	// expost build info
	common.BuildDate = builddate
	glog.Info("go-read version: ", common.BuildDate)

	// sqlite shit
	var db *sql.DB

	dbfileExists := false
	if _, err := os.Stat(*databasefile); err == nil {
		dbfileExists = true
	}

	glog.Info("does dbfile exist? ", dbfileExists)

	if true {
		var err error
		db, err = sql.Open("sqlite3", *databasefile)
		if err != nil {
			glog.Fatal(err)
		}
		defer db.Close()
	}

	if !dbfileExists {
		dao.Init(db)
	}

	// initialise mux router
	r := mux.NewRouter()

	// static files
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/logs/", http.StripPrefix("/logs/", http.FileServer(http.Dir("logs"))))

	// http handles
	r.HandleFunc("/stats/raw", stats_api.Handler)
	r.HandleFunc("/stats", handlers.StatsHandler(db))
	r.HandleFunc("/add", handlers.AddArticle(db))
	r.HandleFunc("/rss", handlers.GenerateRSS(db))
	r.HandleFunc("/", handlers.IndexPage)

	// start server
	http.Handle("/", r)
	glog.Infof("running on port %d", *port)
	err := http.ListenAndServe(":"+strconv.Itoa(*port), log(http.DefaultServeMux))
	if err != nil {
		glog.Fatal(err)
	}
}
