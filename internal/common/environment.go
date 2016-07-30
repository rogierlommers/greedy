package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// BuildDate is the date when the greedy binary has been build
	BuildDate string

	// CommitHash shows the GIT commit hash at build time
	CommitHash string

	// Port where servers runs ons
	Port int

	// Host where server runs on
	Host string

	// Databasefile contains path to storage
	Databasefile     string
	FeedsAuthorName  = "your name"
	FeedsAuthorEmail = "your-email@address.com"
	FeedsLink        = "your-site.com"
)

// ReadEnvironment reads properies form environment. If no environment
// properties are set, then then defaults will be loaded
func ReadEnvironment() {
	// override configuration with environment vars
	// example: GREEDY_PORT=/tmp/greedy.sqlite
	viper.SetEnvPrefix("GREEDY")
	viper.SetDefault("databasefile", "articles.bolt")
	viper.SetDefault("port", 8080)
	viper.SetDefault("host", "0.0.0.0")
	viper.AutomaticEnv()

	Port = viper.GetInt("port")
	Databasefile = viper.GetString("databasefile")
	Host = viper.GetString("host")

	log.Infof("environment loaded [host: %s], [port: %d], [databasefile: %s]", Host, Port, Databasefile)
	log.Infof("greedy info [builddate: %s], [git commit hash: %s]", BuildDate, CommitHash)
}
