package common

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/emicklei/go-selfdiagnose"
)

type ReportBuildAndDate struct {
	CommitHash string
	Date       string
}

func (r ReportBuildAndDate) Run(ctx *selfdiagnose.Context, result *selfdiagnose.Result) {
	result.Passed = true
	result.Reason = fmt.Sprintf("builddate: %s, commithash: %s", r.Date, r.CommitHash)
}

func (r ReportBuildAndDate) Comment() string {
	return "build information"
}

func (r ReportBuildAndDate) Timeout() time.Duration {
	return 0 * time.Second
}

func SetupSelfdiagnose() {
	http.HandleFunc("/selfdiagnose", handleSelfdiagnose)

	// Display commandline args
	selfdiagnose.Register(ArgumentList(os.Args))

	// Inspect the HTTP request
	selfdiagnose.Register(ReportHttpRequest{})

	// What is my hostname
	selfdiagnose.Register(ReportHostname{})

	// report build and date
	report := ReportBuildAndDate{
		Date:       BuildDate,
		CommitHash: CommitHash,
	}

	// finally build report
	selfdiagnose.Register(report)
}

func handleSelfdiagnose(w http.ResponseWriter, r *http.Request) {
	ctx := selfdiagnose.NewContext()
	ctx.Variables["http.request"] = r
	selfdiagnose.DefaultRegistry.RunWithContext(selfdiagnose.HtmlReporter{w}, ctx)
}

type ArgumentList []string

func (r ArgumentList) Run(ctx *selfdiagnose.Context, result *selfdiagnose.Result) {
	result.Passed = true
	var arguments string
	for _, element := range r {
		arguments += element + "<br/>"
	}
	result.Reason = fmt.Sprintf(arguments)
}

func (r ArgumentList) Comment() string {
	return "commandline args"
}

func (r ArgumentList) Timeout() time.Duration {
	return 0 * time.Second
}

type ReportHttpRequest struct{}

func (r ReportHttpRequest) Run(ctx *selfdiagnose.Context, result *selfdiagnose.Result) {
	req, ok := ctx.Variables["http.request"]
	if !ok {
		result.Passed = false
		result.Reason = "missing variable 'http.request'"
		return
	}
	var buf bytes.Buffer
	for k, v := range req.(*http.Request).Header {
		buf.WriteString(fmt.Sprintf("%s = %s<br/>", k, v))
	}
	result.Passed = true
	result.Reason = buf.String()
}

func (r ReportHttpRequest) Comment() string { return "headers of this Http request" }

func (r ReportHttpRequest) Timeout() time.Duration { return 0 * time.Second }

type ReportHostname struct{}

func (r ReportHostname) Comment() string { return "hostname as reported by the os" }

func (r ReportHostname) Timeout() time.Duration { return 0 * time.Second }

func (r ReportHostname) Run(ctx *selfdiagnose.Context, result *selfdiagnose.Result) {
	h, err := os.Hostname()
	result.Passed = err == nil
	result.Reason = h
}
