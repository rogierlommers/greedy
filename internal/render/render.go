package render

import (
	"html/template"
	"net/http"

	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
)

var staticBox *rice.Box

func CreateStaticBox() {

	// create rice.box with static files
	staticBox = rice.MustFindBox("../../static")

	// css files should be exposed through /css endpoint
	cssFileServer := http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox()))
	http.Handle("/static/", cssFileServer)
}

func DisplayPage(w http.ResponseWriter, r *http.Request, dynamicData interface{}) {
	templateString, err := staticBox.String("index.tmpl")
	if err != nil {
		log.Panicf("render error: %s", err)
	}

	tmplMessage, err := template.New("messsage").Parse(templateString)
	if err != nil {
		log.Panicf("render error: %s", err)
	}

	tmplMessage.Execute(w, dynamicData)
}
