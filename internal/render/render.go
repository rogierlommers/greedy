package render

import (
	"log"
	"net/http"

	"html/template"
)

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
