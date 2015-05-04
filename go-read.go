package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/rogierlommers/go-read/internal/common"
	"github.com/rogierlommers/go-read/internal/model"
	"github.com/rogierlommers/go-read/internal/rest"
)

var databasefile = flag.String("databasefile", "database.xml", "XML file where items are stored")
var port = flag.Int("port", 8080, "http listener port")

func init() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	common.DatabaseFile = *databasefile
	model.CreateDatabaseIfNotExists()
}

func main() {
	defer glog.Flush()

	// read database
	database := model.ReadFileIntoSlice()

	// initialise mux router
	r := mux.NewRouter()

	// static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/html"))))

	// http handles
	r.HandleFunc("/stats", rest.StatsHandler(&database))
	r.HandleFunc("/add/{base64url}", rest.AddArticle(&database))
	r.HandleFunc("/rss", rest.GenerateRSS(&database))
	r.HandleFunc("/", rest.IndexPage)

	// start server
	http.Handle("/", r)
	glog.Infof("running on port %d", *port)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
