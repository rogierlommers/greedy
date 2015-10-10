package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/dao"
	"github.com/rogierlommers/greedy/internal/handlers"
	"github.com/rogierlommers/greedy/internal/render"
	log "gopkg.in/inconshreveable/log15.v2"
)

var BuildDate string

func main() {
	// read environment vars
	common.BuildDate = BuildDate
	common.ReadEnvironment()

	// initialize bolt storage
	dao.Open()
	defer dao.Close()

	// initialise mux router
	router := mux.NewRouter()

	// selfdiagnose
	common.SetupSelfdiagnose()

	// setup statics
	render.CreateStaticBox()

	// http handles

	router.HandleFunc("/add", handlers.AddArticle)
	router.HandleFunc("/", handlers.IndexPage)
	//	router.HandleFunc("/stats", handlers.StatsHandler(db))
	//	router.HandleFunc("/export", handlers.ExportCSV(db))

	//	router.HandleFunc("/rss", handlers.GenerateRSS(db))
	//

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
