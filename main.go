package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"github.com/rogierlommers/go-read/internal/common"
	"github.com/rogierlommers/go-read/internal/dao"
)

// TODO
// get hostname from request

// SOURCES
// https://github.com/siadat/eton
// INIT.D SCRIPT HIER: https://github.com/samwierema/go-url-shortener
// https://github.com/samwierema?tab=repositories

// injected by the build process
var builddate = "unknown build date"

// read flags
var databasefile = flag.String("databasefile", "database.xml", "XML file where items are stored")
var port = flag.Int("port", 8080, "http listener port")

func init() {
	flag.Parse()
	flag.Lookup("alsologtostderr").Value.Set("true")
	common.DatabaseFile = *databasefile
	dao.CreateDatabaseIfNotExists()
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

	// read database
	database := dao.ReadFileIntoSlice()

	// initialise mux router
	r := mux.NewRouter()

	// static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/html"))))
	http.Handle("/logs/", http.StripPrefix("/logs/", http.FileServer(http.Dir("logs"))))

	// http handles
	r.HandleFunc("/stats/raw", stats_api.Handler)
	r.HandleFunc("/stats", dao.StatsHandler(&database))
	r.HandleFunc("/add/{base64url}", dao.AddArticle(&database))
	r.HandleFunc("/rss", dao.GenerateRSS(&database))
	r.HandleFunc("/", dao.IndexPage)

	// start server
	http.Handle("/", r)
	glog.Infof("running on port %d", *port)
	err := http.ListenAndServe(":"+strconv.Itoa(*port), log(http.DefaultServeMux))
	if err != nil {
		glog.Fatal(err)
	}
}
