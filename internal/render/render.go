package render

import (
	"log"
	"net/http"

	"github.com/GeertJohan/go.rice"

	"html/template"
)

var staticBox *rice.Box

func CreateStaticBox() {

	// create rice.box with static files
	staticBox = rice.MustFindBox("../../static")

	// css files should be exposed through /css endpoint
	cssFileServer := http.StripPrefix("/css/", http.FileServer(staticBox.HTTPBox()))
	http.Handle("/css/", cssFileServer)
}

func DisplayPage(w http.ResponseWriter, r *http.Request, dynamicData interface{}) {
	templateString, err := staticBox.String("index.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	tmplMessage, err := template.New("messsage").Parse(templateString)
	if err != nil {
		log.Fatal(err)
	}

	tmplMessage.Execute(w, dynamicData)
}
