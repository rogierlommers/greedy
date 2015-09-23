package selfdiagnose

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import "log"

// LoggingReporter reports on the task by printing using standard log.
type LoggingReporter struct{}

func (l LoggingReporter) Report(results []*Result) {
	for _, each := range results {
		verdict := "passed"
		if !each.Passed {
			verdict = "failed"
		}
		comment, separator := "", ""
		if len(each.Target.Comment()) > 0 {
			comment = each.Target.Comment()
			separator = "."
		}
		log.Printf("[%s] %s%s %s", verdict, comment, separator, each.Reason)
	}
}
