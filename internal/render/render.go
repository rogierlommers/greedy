package render

import (
	"net/http"
	"path"
	"text/template"
)

func DisplayPage(w http.ResponseWriter, r *http.Request, text string, page string) {

	// determine page
	// show page

	fp := path.Join("static", "templates", "stats.html")
	tmpl, parseErr := template.ParseFiles(fp)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusInternalServerError)
		return
	}

	obj := map[string]string{"statsmessage": stats}

	if templErr := tmpl.Execute(w, obj); templErr != nil {
		http.Error(w, templErr.Error(), http.StatusInternalServerError)
	}

}
