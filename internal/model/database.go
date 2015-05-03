package model

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/rogierlommers/go-read/internal/common"
)

func CreateDatabaseIfNotExists() {
	glog.Info("database: ", common.DatabaseFile)
	if _, err := os.Stat(common.DatabaseFile); os.IsNotExist(err) {
		glog.Infof("creating file: %s", common.DatabaseFile)
		_, err := os.Create(common.DatabaseFile)
		if err != nil {
			glog.Errorf("error creating database file -> '%s'", err)
			os.Exit(1)
		}
	}
}

func ReadFileIntoSlice() ReadingListRecords {
	xmlFile, err := os.Open(common.DatabaseFile)
	defer xmlFile.Close()

	if err != nil {
		glog.Errorf("error reading records from file -> %s", err)
		os.Exit(1)
	}

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var database ReadingListRecords
	xml.Unmarshal(XMLdata, &database)

	glog.Infof("read %d records in struct", len(database.Records))
	common.LastID = len(database.Records)
	return database
}
