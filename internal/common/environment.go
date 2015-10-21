package common

import (
	"github.com/spf13/viper"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	BuildDate        string
	CommitHash       string
	Port             int
	Host             string
	Databasefile     string
	FeedsAuthorName  = "Rogier Lommers"
	FeedsAuthorEmail = "rogier@lommers.org"
	FeedsLink        = "http://www.lommers.org"
)

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
