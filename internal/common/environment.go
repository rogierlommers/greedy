package common

import (
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

var (
	BuildDate        string
	Port             int
	Host             string
	Databasefile     string
	FeedsAuthorName  = "Rogier Lommers"
	FeedsAuthorEmail = "rogier@lommers.org"
	FeedsLink        = "http://www.lommers.org"
)

func ReadEnvironment() {
	viper.SetEnvPrefix("GREEDY")
	viper.SetDefault("databasefile", "articles.sqlite")
	viper.SetDefault("port", 8080)
	viper.SetDefault("host", "0.0.0.0")
	viper.AutomaticEnv()

	Port = viper.GetInt("port")
	Databasefile = viper.GetString("databasefile")
	Host = viper.GetString("host")

	glog.Infof("environment var loaded [port=%d]", Port)
	glog.Infof("environment var loaded [databasefile=%s]", Databasefile)
	glog.Infof("environment var loaded [host=%s]", Host)
	glog.Info("greedy version: ", BuildDate)
}
