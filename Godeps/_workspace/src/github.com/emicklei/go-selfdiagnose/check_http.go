package selfdiagnose

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// CheckHttp send a http.Request and check the status code. 200 OK = Passed
type CheckHttp struct {
	BasicTask
	Request      *http.Request
	ShowResponse bool
}

// Run sends the request and updates the result.
func (c CheckHttp) Run(ctx *Context, result *Result) {
	client := new(http.Client)
	resp, err := client.Do(c.Request)
	if err != nil {
		result.Passed = false
		result.Reason = err.Error()
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		result.Passed = false
		result.Reason = resp.Status
		return
	}

	summary := fmt.Sprintf("%s %s => %s", c.Request.Method, c.Request.URL.String(), resp.Status)
	result.Passed = true
	if c.ShowResponse {
		var buf bytes.Buffer
		buf.WriteString(summary)
		buf.WriteString("\n\n")
		io.Copy(&buf, resp.Body)
		summary = buf.String()
	}
	result.Reason = summary
}
