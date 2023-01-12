package common

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	BuildVersion     string
	CommitHash       string
	Port             int
	Host             string
	Databasefile     string
	FeedsAuthorName  = "Greedy | your personal reading list"
	FeedsAuthorEmail = ""
	FeedsLink        = ""
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

	logrus.Infof("environment loaded [host: %s], [port: %d], [databasefile: %s]", Host, Port, Databasefile)
	logrus.Infof("greedy info [builddate: %s], [git commit hash: %s]", BuildVersion, CommitHash)

	if os.Getenv("DEV") == "true" {
		logrus.Info("DEVELOPMENT_MODE")
		BuildVersion = "development"
	}
}
