package selfdiagnose

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import "time"

var DefaultRegistry = Registry{}

// Registry holds the collection or registered Tasks. It can run them all.
type Registry struct {
	tasks []Task
}

// Register adds a task to the collection.
func (r *Registry) Register(t Task) {
	r.tasks = append(r.tasks, t)
}

// Run executes all registered task (in order) and reports using a Reporter.
func (r Registry) Run(rep Reporter) {
	r.RunWithContext(rep, NewContext())
}

// RunTasks runs all registered tasks and returns the result of each.
func (r Registry) RunTasks(ctx *Context) []*Result {
	results := []*Result{}
	for _, each := range r.tasks {
		resultCh := make(chan *Result, 1)
		now := time.Now()
		go func() {
			defer close(resultCh)
			res := new(Result)
			if t, ok := each.(HasSeverity); ok {
				res.Severity = t.Severity()
			}
			res.Target = each
			each.Run(ctx, res)
			resultCh <- res
		}()
		timeout := 1 * time.Second
		// task can override the non-zero value
		if t, ok := each.(HasTimeout); ok {
			if t.Timeout() > 0 {
				timeout = t.Timeout()
			}
		}
		var result *Result
		select {
		case <-time.After(timeout):
			// task took longer than timeout so create result to report that
			res := new(Result)
			if t, ok := each.(HasSeverity); ok {
				res.Severity = t.Severity()
			}
			res.Target = each
			res.Passed = false
			res.Reason = "task did not complete within timeout"
			result = res
		case result, _ = <-resultCh:
		}
		result.CompletedIn = time.Now().Sub(now)
		results = append(results, result)
	}
	return results
}

func (r Registry) RunWithContext(rep Reporter, ctx *Context) {
	rep.Report(r.RunTasks(ctx))
}

// Register adds a task to the default registry
func Register(t Task) {
	DefaultRegistry.Register(t)
}

// Run delegates to the DefaultRegistry
func Run(rep Reporter) {
	DefaultRegistry.Run(rep)
}
