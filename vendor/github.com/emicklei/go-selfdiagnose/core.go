package selfdiagnose

// Copyright 2015-2016 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import "time"

// VERSION is used for including it in a report.
const VERSION = "1.6.7"

var since = time.Now()

// Task describes a diagnostic task that can be run.
type Task interface {
	Run(ctx *Context, result *Result)
	Comment() string
}

// HasTimeout indicates whether a task can be timedout.
type HasTimeout interface {
	Timeout() time.Duration
}

// Result captures the execution result of a Task.
type Result struct {
	Target Task
	Passed bool

	// Reason can be a string, template.HTML or any printable value.
	Reason      interface{}
	CompletedIn time.Duration
	Severity    Severity
}

// Context can be used to read/write variable during the execution of a selfdiagnose run.
type Context struct {
	Variables map[string]interface{}
}

// Severity indicates whether a failing task is severe or not. See constants for possible values.
type Severity string

// HasSeverity indicates whether a task has a severity level.
type HasSeverity interface {
	Severity() Severity
}

const (
	// SeverityNone indicates that a failed result is not severe
	SeverityNone Severity = "none"
	// SeverityWarning indicates that a failed result is some to look it but not urgent.
	SeverityWarning Severity = "warning"
	// SeverityCritical indicates that a failed result is making the system unhealthy.
	SeverityCritical Severity = "critical"
)

// BasicTask can be used to embed in new created Task types.
type BasicTask struct {
	comment  string
	timeout  time.Duration
	severity Severity
}

// Comment returns a short description what this task is for.
func (t BasicTask) Comment() string {
	return t.comment
}

// SetComment is to set the comment
func (t *BasicTask) SetComment(text string) {
	t.comment = text
}

func (t BasicTask) Timeout() time.Duration {
	return t.timeout
}

func (t *BasicTask) SetTimeout(after time.Duration) {
	t.timeout = after
}

func (t *BasicTask) SetSeverity(s Severity) {
	t.severity = s
}

func (t BasicTask) Severity() Severity {
	if len(t.severity) == 0 {
		return SeverityCritical
	}
	return t.severity
}

// NewContext creates a new empty Context to run tasks.
func NewContext() *Context {
	return &Context{map[string]interface{}{}}
}

// Reporter describes how to report task execution results.
type Reporter interface {
	Report(results []*Result)
}
