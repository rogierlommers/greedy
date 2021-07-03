package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	BuildDate        string
	CommitHash       string
	Port             int
	Host             string
	Databasefile     string
	ToEmail          string
	FromEmail        string
	SMTPHost         string
	SMTPUser         string
	SMTPPassword     string
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
	ToEmail = viper.GetString("to_email")
	FromEmail = viper.GetString("from_email")
	SMTPHost = viper.GetString("smtp_host")
	SMTPUser = viper.GetString("smtp_user")
	SMTPPassword = viper.GetString("smtp_password")

	Host = viper.GetString("host")

	log.Infof("environment loaded [host: %s], [port: %d], [databasefile: %s]", Host, Port, Databasefile)
	log.Infof("greedy info [builddate: %s], [git commit hash: %s]", BuildDate, CommitHash)
}
