package common

import (
	"github.com/spf13/viper"
	log "gopkg.in/inconshreveable/log15.v2"
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

	Databasefile     string
	FeedsAuthorName  = "Rogier Lommers"
	FeedsAuthorEmail = "rogier@lommers.org"
	FeedsLink        = "http://www.lommers.org"
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

	log.Info("environment", "host", Host, "port", Port, "databasefile", Databasefile)
	log.Info("greedy", "builddate", BuildDate, "commithash", CommitHash)
}
