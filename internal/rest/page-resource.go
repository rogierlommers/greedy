package rest

import (
	"html/template"
	"net/http"
	"path"

	"github.com/golang/glog"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	glog.Info("opening stats page")
}

func RSS(w http.ResponseWriter, r *http.Request) {
	glog.Info("opening RSS page")
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	glog.Info("opening index page")
	glog.Info(r)
	fp := path.Join("static", "templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj := map[string]string{"serverLocation": "http://localhost:8080"}

	if err := tmpl.Execute(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
