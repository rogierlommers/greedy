package selfdiagnose

// Copyright 2016 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"fmt"
	"strconv"
	"time"
)

type runReport struct {
	SelfDiagnose map[string]interface{} `json:"selfdiagnose" `
	Run          time.Time              `json:"run" `
	Results      []taskReport           `json:"results" `
}

type taskReport struct {
	Task     string      `json:"task" xml:"task"`
	Status   string      `json:"status" xml:"status"`
	Comment  string      `json:"comment" xml:"comment"`
	Message  interface{} `json:"message" xml:"message"`
	Duration string      `json:"duration" xml:"duration"`
	Severity string      `json:"severity" xml:"severity"`
}

func buildTaskReports(results []*Result) (list []taskReport) {
	for _, each := range results {
		list = append(list, taskReport{
			Task:     fmt.Sprintf("%T", each.Target),
			Status:   statusLabel(each.Passed),
			Comment:  each.Target.Comment(),
			Message:  each.Reason,
			Duration: toMillisecondsString(each.CompletedIn),
			Severity: string(each.Severity),
		})
	}
	return
}

func buildRunReport(results []*Result) runReport {
	c, f := checksAndFailures(results)
	return runReport{
		SelfDiagnose: map[string]interface{}{
			"version":     VERSION,
			"since":       since.String(),
			"completedIn": toMillisecondsString(totalDuration(results)),
			"checks":      c,
			"failures":    f,
		},
		Run:     time.Now(),
		Results: buildTaskReports(results),
	}
}

func statusLabel(ok bool) string {
	if !ok {
		return "FAILED"
	}
	return "OK"
}

func totalDuration(results []*Result) (total time.Duration) {
	for _, each := range results {
		total += each.CompletedIn
	}
	return total
}

func toMillisecondsString(d time.Duration) string {
	return strconv.FormatInt(d.Nanoseconds()/1000000, 10) // ms
}

func checksAndFailures(results []*Result) (checks, failures int) {
	for _, r := range results {
		checks++
		if !r.Passed {
			failures++
		}
	}
	return
}
