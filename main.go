package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rogierlommers/greedy/internal/articles"
	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/render"
	log "gopkg.in/inconshreveable/log15.v2"
)

var BuildDate string

// http://bl.ocks.org/joyrexus/22c3ef0984ed957f54b9
// Skybox Analytics - https://github.com/skybox/skybox
// DVID - Added Bolt as optional storage engine and testing it against Basho-tuned leveldb. https://github.com/janelia-flyem/dvid

func main() {
	// read environment vars
	common.BuildDate = BuildDate
	common.ReadEnvironment()

	// initialize bolt storage
	articles.Open()
	defer articles.Close()

	// initialise mux router
	router := mux.NewRouter()

	// selfdiagnose
	common.SetupSelfdiagnose()

	// setup statics
	render.CreateStaticBox()

	// http handles
	router.HandleFunc("/", articles.IndexPage)
	router.HandleFunc("/add", articles.AddArticle)
	router.HandleFunc("/rss", articles.DisplayRSS)

	//router.HandleFunc("/stats", handlers.StatsHandler)
	//	router.HandleFunc("/export", handlers.ExportCSV(db))

	// start cleanup db routing
	//	go dao.Cleanup(db)

	// start server
	http.Handle("/", router)
	log.Info("deamon listening", "host", common.Host, "port", common.Port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", common.Host, common.Port), nil)
	if err != nil {
		log.Crit("daemon could not bind on interface", "host", common.Host, "port", common.Port)
	}
}
