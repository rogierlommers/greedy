package main

import (
	"flag"
	"net/http"

	"fmt"

	"github.com/GeertJohan/go.rice"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/dao"
	"github.com/rogierlommers/greedy/internal/handlers"
)

func init() {
	flag.Parse()
	flag.Lookup("alsologtostderr").Value.Set("true")
}

var BuildDate string

func main() {
	defer glog.Flush()

	// read environment vars
	common.BuildDate = BuildDate
	common.ReadEnvironment()

	// initialize sqlite storage
	db := dao.Init(common.Databasefile)
	defer db.Close()

	// initialise mux router
	router := mux.NewRouter()

	// selfdiagnose
	common.SetupSelfdiagnose()

	// static files
	cssBox := rice.MustFindBox("static-css")
	cssFileServer := http.StripPrefix("/css/", http.FileServer(cssBox.HTTPBox()))
	http.Handle("/css/", cssFileServer)

	// http handles
	router.HandleFunc("/stats", handlers.StatsHandler(db))
	router.HandleFunc("/export", handlers.ExportCSV(db))
	router.HandleFunc("/add", handlers.AddArticle(db))
	router.HandleFunc("/rss", handlers.GenerateRSS(db))
	router.HandleFunc("/", handlers.IndexPage)

	// start cleanup db routing
	go dao.Cleanup(db)

	// start server
	http.Handle("/", router)
	hostPort := fmt.Sprintf("%s:%d", common.Host, common.Port)
	glog.Infof("running on: %s", hostPort)

	err := http.ListenAndServe(hostPort, nil)
	if err != nil {
		glog.Fatal(err)
	}
}
