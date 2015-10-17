package articles

import (
	"net/http"
	"time"

	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/render"
	log "gopkg.in/inconshreveable/log15.v2"
)

func AddArticle(w http.ResponseWriter, r *http.Request) {
	queryParam := r.FormValue("url")
	if len(queryParam) == 0 || queryParam == "about:blank" {

		renderObject := map[string]interface{}{
			"IsErrorPage":  "true",
			"errorMessage": "unable to insert empty or about:blank page",
		}
		render.DisplayPage(w, r, renderObject)
		return
	}

	newArticle := Article{
		//ID:    getMD5Hash(queryParam), // function will add time.Now() to make it unique
		Url:   queryParam,
		Added: time.Now(),
	}

	err := newArticle.Save()
	if err != nil {
		log.Warn("error saving article", "hostname", getHostnameFromUrl(queryParam), "id", newArticle.ID)
	}

	// finally output confirmation page
	renderObject := map[string]interface{}{
		"IsConfirmation": "true",
		"hostname":       getHostnameFromUrl(queryParam),
	}
	render.DisplayPage(w, r, renderObject)
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	renderObject := map[string]interface{}{
		"IsLandingPage":  "true",
		"serverLocation": r.Host,
		"buildversion":   common.BuildDate,
	}
	render.DisplayPage(w, r, renderObject)
}
