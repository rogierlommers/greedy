package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"github.com/rogierlommers/go-read/internal/common"
	"github.com/rogierlommers/go-read/internal/dao"
)

// TODO
// get hostname from request

// SOURCES
// http://stackoverflow.com/questions/30037515/adding-element-to-slice-in-handlerfunc-and-return-as-a-whole

// injected by the build process
var version = "unknown version"
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

func main() {
	defer glog.Flush()

	// expost build info
	common.BuildVersion = version
	common.BuildDate = builddate

	// read database
	database := dao.ReadFileIntoSlice()

	// initialise mux router
	r := mux.NewRouter()

	// static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/html"))))

	// http handles
	r.HandleFunc("/stats", dao.StatsHandler(&database))
	r.HandleFunc("/add/{base64url}", dao.AddArticle(&database))
	r.HandleFunc("/rss", dao.GenerateRSS(&database))
	r.HandleFunc("/", dao.IndexPage)

	// start server
	http.Handle("/", r)
	glog.Infof("running on port %d", *port)
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		glog.Fatal(err)
	}
}
