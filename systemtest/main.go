package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	log "gopkg.in/inconshreveable/log15.v2"
)

var amount = 100

func serveTestPage() {
	log.Info("systemtest", "status", "start serving pages")
	http.HandleFunc("/", someHandler)
	http.ListenAndServe(":8081", nil)
}

func someHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	bodyCount := r.FormValue("bodyCount")
	htmlTemplate := template.Must(template.New("Page").Parse(`<html><head><title>{{.title}}</title></head><body>body</body></html>`))
	renderObject := map[string]interface{}{
		"title":   title,
		"counter": bodyCount,
	}
	htmlTemplate.Execute(w, renderObject)
}

func main() {
	// run goroute for serving pages
	go serveTestPage()

	var url string
	for i := 0; i < amount; i++ {
		url = fmt.Sprintf("http://localhost:8080/add?url=http://localhost:8081?title=%s", strconv.Itoa(i))
		response, err := http.Get(url)
		defer response.Body.Close()
		log.Info("systemtest", "url", url, "response", response.Status, "message", err)
	}

	log.Info("systemtest", "status", "amount", amount)
}
