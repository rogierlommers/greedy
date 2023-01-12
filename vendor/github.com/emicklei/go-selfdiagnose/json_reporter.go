package selfdiagnose

// Copyright 2016 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"encoding/json"
	"io"
)

// JSONReporter is to produce a JSON report and it written on an io.Writer.
type JSONReporter struct {
	Writer io.Writer
}

// Report produces a JSON report including a summary
func (j JSONReporter) Report(results []*Result) {
	// silently ignore the errors
	data, _ := json.MarshalIndent(buildRunReport(results), "", "\t")
	j.Writer.Write(data)
}
