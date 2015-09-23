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

func (r Registry) RunWithContext(rep Reporter, ctx *Context) {
	results := []*Result{}
	for _, each := range r.tasks {
		resultCh := make(chan *Result, 1)
		now := time.Now()
		go func() {
			res := new(Result)
			res.Target = each
			each.Run(ctx, res)
			resultCh <- res // will not block if closed
		}()
		timeout := each.Timeout()
		if timeout == 0 {
			timeout = 1 * time.Second
		}
		var result *Result
		select {
		case <-time.After(timeout):
			res := new(Result)
			res.Target = each
			res.Passed = false
			res.Reason = "task did not complete within timeout"
			result = res
		case result, _ = <-resultCh:
		}
		result.CompletedIn = time.Now().Sub(now)
		results = append(results, result)
	}
	rep.Report(results)
}

// Register adds a task to the default registry
func Register(t Task) {
	DefaultRegistry.Register(t)
}

// Run delegates to the DefaultRegistry
func Run(rep Reporter) {
	DefaultRegistry.Run(rep)
}
