package render

import (
	"net/http"
	"path"
	"text/template"
)

func DisplayPage(w http.ResponseWriter, r *http.Request, renderObject map[string]string, templateFile string) {
	fp := path.Join("static", "templates", templateFile)
	tmpl, parseErr := template.ParseFiles(fp)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusInternalServerError)
		return
	}

	if templErr := tmpl.Execute(w, renderObject); templErr != nil {
		http.Error(w, templErr.Error(), http.StatusInternalServerError)
	}
}
