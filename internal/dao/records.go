package dao

import (
	"encoding/xml"
	"errors"
	"time"

	"github.com/golang/glog"
)

type Record struct {
	XMLName xml.Name `xml:"record"`
	Id      int      `xml:"id"`
	URL     string   `xml:"url"`
}

type ReadingListRecords struct {
	XMLName xml.Name `xml:"records"`
	Records []Record `xml:"record"`
}

type RecordSet interface {
	Append(record Record) error
}

func (records *ReadingListRecords) Append(record Record) error {
	newRecords := append(records.Records, record)

	if newRecords == nil {
		return errors.New("Could not append record")
	} else {
		records.Records = newRecords
		return nil
	}
}

// ById implements sort.Interface for []Records based on the Id field
type ById []Record

func (a ById) Len() int           { return len(a) }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }

func addAdditionalInfo() {
	// https://github.com/moovweb/gokogiri
	// https://github.com/PuerkitoBio/goquery
	glog.Info("starting thread")
	time.Sleep(10 * time.Second)
	glog.Info("ending thread")
}
