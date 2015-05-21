package render

import (
	"net/http"
	"path"
	"text/template"
)

func DisplayPage(w http.ResponseWriter, r *http.Request, renderObject map[string]string, templateFile string) {
	templateDirectory := path.Join("static", "templates")
	baseFile := path.Join(templateDirectory, templateFile)
	headerFile := path.Join(templateDirectory, "_header.html")
	footerFile := path.Join(templateDirectory, "_footer.html")

	tmpl, parseErr := template.ParseFiles(baseFile, headerFile, footerFile)
	if parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusInternalServerError)
		return
	}

	if templErr := tmpl.Execute(w, renderObject); templErr != nil {
		http.Error(w, templErr.Error(), http.StatusInternalServerError)
	}
}
