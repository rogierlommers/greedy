package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/rogierlommers/go-read/internal/rest"
)

var databasefile = flag.String("databasefile", "database.xml", "XML file where items are stored")
var port = flag.Int("port", 8080, "http listener port")

func main() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	defer glog.Flush()

	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// selfdiagnose
	http.HandleFunc("/stats/", rest.StatsHandler)

	// rss output
	http.HandleFunc("/rss/", rest.RSS)

	// serve static files
	http.HandleFunc("/", rest.IndexPage)

	// start server
	glog.Infof("running on port %d", *port)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
