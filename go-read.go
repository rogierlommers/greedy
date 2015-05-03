package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
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

// http://www.gorillatoolkit.org/pkg/mux

func main() {
	defer glog.Flush()

	// read database
	database := model.ReadFileIntoSlice()
	spew.Dump(database)

	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// selfdiagnose
	http.HandleFunc("/stats/", rest.StatsHandler)

	// rss output
	http.HandleFunc("/rss/", rest.RSS)

	// add new url
	http.HandleFunc("/add/:url", rest.AddArticle)

	// serve static files
	http.HandleFunc("/index", rest.IndexPage)

	// start server
	glog.Infof("running on port %d", *port)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
