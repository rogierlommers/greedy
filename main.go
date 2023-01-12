package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rogierlommers/greedy/internal/articles"
	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/render"
	"github.com/sirupsen/logrus"
)

func main() {
	// read environment vars and setup http client
	common.ReadEnvironment()

	// initialize bolt storage
	articles.Open()
	defer articles.Close()

	// initialize mux router
	router := mux.NewRouter()

	// setup statics
	render.CreateStaticBox()

	// http handles
	router.HandleFunc("/", articles.IndexPage)
	router.HandleFunc("/add", articles.AddArticle)
	router.HandleFunc("/rss", articles.DisplayRSS)
	router.HandleFunc("/export", articles.ExportCSV)

	// schedule cleanup routing
	articles.ScheduleCleanup()

	// start server
	http.Handle("/", router)
	logrus.Infof("deamon running on host %s and port %d", common.Host, common.Port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", common.Host, common.Port), nil)
	if err != nil {
		logrus.Panicf("daemon could not bind on interface: %s, port: %d", common.Host, common.Port)
	}
}
