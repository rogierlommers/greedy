package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rogierlommers/greedy/internal/articles"
	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/render"
)

func main() {
	// read environment vars and setup http client
	common.ReadEnvironment()
	articles.NewClient()

	// initialize bolt storage
	articles.Open()
	defer articles.Close()

	// initialize mux router
	router := mux.NewRouter()

	// selfdiagnose
	common.SetupSelfdiagnose()

	// setup statics
	render.CreateStaticBox()

	// http handles
	router.HandleFunc("/", articles.IndexPage)
	router.HandleFunc("/add", articles.AddArticle)
	router.HandleFunc("/rss", articles.DisplayRSS)
	router.HandleFunc("/export", articles.ExportCSV)

	// start server
	http.Handle("/", router)
	log.Infof("deamon running on host %s and port %d", common.Host, common.Port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", common.Host, common.Port), nil)
	if err != nil {
		log.Panicf("daemon could not bind on interface: %s, port: %d", common.Host, common.Port)
	}
}
