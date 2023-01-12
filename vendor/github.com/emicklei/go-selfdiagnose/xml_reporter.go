package selfdiagnose

// Copyright 2016 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"encoding/xml"
	"io"
	"time"
)

// XMLReporter is to produce a XML report and it written on an io.Writer.
type XMLReporter struct {
	Writer io.Writer
}

// Report produces a XML report including a summary
func (x XMLReporter) Report(results []*Result) {
	c, f := checksAndFailures(results)
	// silently ignore the errors
	r := xmlReport{
		Run:         time.Now(),
		Since:       since,
		Version:     VERSION,
		CompletedIn: toMillisecondsString(totalDuration(results)),
		Checks:      c,
		Failures:    f,
		Results:     buildTaskReports(results),
	}
	io.WriteString(x.Writer, xml.Header)
	data, _ := xml.MarshalIndent(r, "  ", "    ")
	x.Writer.Write(data)
}

type xmlReport struct {
	XMLName     xml.Name     `json:"-" xml:"selfdiagnose"`
	Run         time.Time    `xml:"run,attr"`
	Since       time.Time    `xml:"since,attr"`
	Version     string       `xml:"version,attr"`
	CompletedIn string       `xml:"completedIn,attr"`
	Checks      int          `xml:"checks,attr"`
	Failures    int          `xml:"failures,attr"`
	Results     []taskReport `xml:"results"`
}
