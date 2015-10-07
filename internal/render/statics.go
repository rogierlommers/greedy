package render

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
)

func RegisterCSS() {
	cssBox := rice.MustFindBox("../../static-css")
	cssFileServer := http.StripPrefix("/css/", http.FileServer(cssBox.HTTPBox()))
	http.Handle("/css/", cssFileServer)

	//	templateBox := rice.MustFindBox("static-templates")
	//	templateFileServer := http.StripPrefix("/css/", http.FileServer(cssBox.HTTPBox()))

	// find a rice.Box
	//	templateBox, err := rice.FindBox("example-templates")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	// get file contents as string
	//	templateString, err := templateBox.String("message.tmpl")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	// parse and execute the template
	//	tmplMessage, err := template.New("message").Parse(templateString)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	tmplMessage.Execute(os.Stdout, map[string]string{"Message": "Hello, world!"})

}
