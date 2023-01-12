package selfdiagnose

// Copyright 2016 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"net/http"
	"strconv"
	"strings"
)

func AddInternalHandlers() {
	http.HandleFunc("/internal/selfdiagnose.html", HandleSelfdiagnose)
	http.HandleFunc("/internal/selfdiagnose.xml", HandleSelfdiagnose)
	http.HandleFunc("/internal/selfdiagnose.json", HandleSelfdiagnose)
}

// HandleSelfdiagnose runs all registered tasks and reports a HTML,JSON or XML report depending on the requested format.
func HandleSelfdiagnose(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext()
	// prepare for ReportHttpRequest
	ctx.Variables["http.request"] = r
	var reporter Reporter
	if strings.HasSuffix(r.URL.Path, ".json") || r.URL.Query().Get("format") == "json" {
		w.Header().Set("Content-Type", "application/json")
		reporter = JSONReporter{w}
	} else if strings.HasSuffix(r.URL.Path, ".xml") || r.URL.Query().Get("format") == "xml" {
		w.Header().Set("Content-Type", "application/xml")
		reporter = XMLReporter{w}
	} else {
		w.Header().Set("Content-Type", "text/html")
		reporter = HtmlReporter{w}
	}
	results := DefaultRegistry.RunTasks(ctx)
	ok := true
	for _, each := range results {
		ok = ok && each.Passed
	}
	// write header first
	w.Header().Set("X-SelfDiagnose-OK", strconv.FormatBool(ok))
	reporter.Report(results)
}
