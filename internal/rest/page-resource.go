package rest

import (
	"html/template"
	"net/http"
	"path"

	"github.com/golang/glog"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	glog.Info("stats page")
}

func RSS(w http.ResponseWriter, r *http.Request) {
	glog.Info("RSS page")
}

func AddArticle(w http.ResponseWriter, r *http.Request) {
	glog.Info("add article")
	//    r.GET("/add/:url", func(c *gin.Context) {
	//        encodedUrl := c.Params.ByName("url")

	//        urlByteArray, err := base64.StdEncoding.DecodeString(encodedUrl)
	//        if err != nil {
	//            glog.Errorf("error decoding url -> %s", err)
	//            return
	//        }

	//        url := string(urlByteArray[:])
	//        database = AddRecord(database, url)
	//        glog.Infof("amount in db: %d", len(database.Records))
	//        c.String(200, url)
	//    })

}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	glog.Info("index page")
	glog.Info(r)
	fp := path.Join("static", "templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO extract serverlocation from header
	obj := map[string]string{"serverLocation": "http://localhost:8080"}

	if err := tmpl.Execute(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
