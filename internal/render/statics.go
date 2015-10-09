package render

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
)

var staticBox *rice.Box

func CreateStaticBox() {

	// create rice.box with static files
	staticBox = rice.MustFindBox("../../static")

	// css files should be exposed through /css endpoint
	cssFileServer := http.StripPrefix("/css/", http.FileServer(staticBox.HTTPBox()))
	http.Handle("/css/", cssFileServer)
}
