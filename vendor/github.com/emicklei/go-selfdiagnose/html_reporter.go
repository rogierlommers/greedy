package selfdiagnose

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"fmt"
	"html/template"
	"io"
	"time"
)

// HtmlReporter is to produce a HTML report and it written on an io.Writer.
type HtmlReporter struct {
	Writer io.Writer
}

type resultRow struct {
	Comment          string
	Description      template.HTML
	Passed           bool
	Severity         string
	Time             time.Duration
	RowStyle         string
	DescriptionStyle string
}

type resultTable struct {
	Rows        []resultRow
	PassedCount int
	FailedCount int
	CompletedIn time.Duration
	Version     string
	Since       time.Time
	ReportDate  time.Time
}

func (r resultTable) TotalCount() int {
	return r.PassedCount + r.FailedCount
}

func (r resultTable) CSS() template.CSS {
	return template.CSS(css())
}

// Report produces a HTML report including a summary
func (h HtmlReporter) Report(results []*Result) {
	resultTable := buildResultTable(results)
	htmlTemplate.Execute(h.Writer, resultTable)
}

var htmlTemplate = template.Must(template.New("Page").Parse(`
<!DOCTYPE html>
<html lang="en-GB">
	<head>
		<meta charset="UTF-8">
	</head>
<body>
	<style>
		{{.CSS}}	
	</style>
	<table>
		<tr class="odd">
			<th>Comment</th>
			<th>Result</th>
			<th>Message</th>
			<th>Severity</th>
			<th>Time</th>
		</tr>
		{{range .Rows}}
		<tr class="{{.RowStyle}}">
			<td>{{.Comment}}</td>
			<td>{{if .Passed}} OK {{ else }} FAIL {{ end }}</td>	
			<td>{{.Description}}</td>	
			<td>{{.Severity}}</td>	
			<td>{{.Time}}</td>	
		</tr>		
		{{end}}
	</table>
	
	<h4>
		Checks: {{.TotalCount}}, Failures: {{.FailedCount}}, Time: {{.CompletedIn}} |
		go-selfdiagnose {{.Version}} | <a href="?format=xml">XML</a>| <a href="?format=json">JSON</a> | since:{{.Since.Format "2006-01-02 3:04"}} report:{{.ReportDate.Format "2006-01-02 3:04"}} </td>
	</h4>
</body>
</html>`))

func buildResultTable(results []*Result) resultTable {
	rows := []resultRow{}
	passedCount := 0
	failedCount := 0
	completedIn := time.Duration(0)
	for i, each := range results {
		row := resultRow{}
		// make reason a HTML value
		switch each.Reason.(type) {
		case template.HTML:
			row.Description = each.Reason.(template.HTML)
		case string:
			row.Description = template.HTML(each.Reason.(string))
		default:
			row.Description = template.HTML(fmt.Sprintf("%v", each.Reason))
		}
		row.Comment = each.Target.Comment()
		row.Time = each.CompletedIn
		row.Passed = each.Passed
		if i%2 == 0 {
			row.RowStyle = "even "
		} else {
			row.RowStyle = "odd "
		}
		if !each.Passed {
			row.RowStyle += "failed "

			if len(each.Severity) > 0 {
				row.RowStyle += string(each.Severity)
				row.Severity = string(each.Severity)
			}
		}
		if each.Passed {
			row.DescriptionStyle = "passed"
			passedCount++
		} else {
			row.DescriptionStyle = "failed"
			failedCount++
		}
		rows = append(rows, row)
		completedIn += each.CompletedIn
	}

	resultTable := resultTable{
		Rows:        rows,
		PassedCount: passedCount,
		FailedCount: failedCount,
		CompletedIn: completedIn,
		Version:     VERSION,
		Since:       since,
		ReportDate:  time.Now(),
	}
	return resultTable

}
