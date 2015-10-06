package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/dao"
	"github.com/rogierlommers/greedy/internal/handlers"
)

// injected by the build process
var BuildDate = "unknown build date"

// read flags
var (
	databasefile = flag.String("databasefile", "articles.db", "sqlite file where items are stored")
	port         = flag.Int("port", 8080, "http listener port")
)

func init() {
	flag.Parse()
	flag.Lookup("alsologtostderr").Value.Set("true")
}

func main() {
	defer glog.Flush()

	// expose build info
	common.BuildDate = BuildDate
	glog.Info("go-read version: ", common.BuildDate)

	// initialize sqlite storage
	db := dao.Init(*databasefile)
	defer db.Close()

	// initialise mux router
	r := mux.NewRouter()

	// selfdiagnose
	common.SetupSelfdiagnose()

	// static files
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))
	http.Handle("/logs/", http.StripPrefix("/logs/", http.FileServer(http.Dir("logs"))))

	// http handles
	r.HandleFunc("/stats", handlers.StatsHandler(db))
	r.HandleFunc("/export", handlers.ExportCSV(db))
	r.HandleFunc("/add", handlers.AddArticle(db))
	r.HandleFunc("/rss", handlers.GenerateRSS(db))
	r.HandleFunc("/", handlers.IndexPage)

	// start cleanup db routing
	go dao.Cleanup(db)

	// start server
	http.Handle("/", r)
	glog.Infof("running on port %d", *port)
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		glog.Fatal(err)
	}
}
