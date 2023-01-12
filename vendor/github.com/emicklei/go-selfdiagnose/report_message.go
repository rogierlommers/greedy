package selfdiagnose

// ReportMessage simply does what it says. Can have comment too.
type ReportMessage struct {
	BasicTask
	Message string
}

func (r ReportMessage) Run(ctx *Context, result *Result) {
	result.Passed = true
	result.Reason = r.Message
	result.Severity = SeverityNone
}
