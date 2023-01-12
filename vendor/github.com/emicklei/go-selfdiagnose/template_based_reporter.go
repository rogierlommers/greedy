package selfdiagnose

import (
	"io"
	"text/template"
)

// TemplateBasedReporter produces a report and it written on an io.Writer.
type TemplateBasedReporter struct {
	Writer   io.Writer
	Template *template.Template
}

// Report produces a report including a summary
func (h TemplateBasedReporter) Report(results []*Result) {
	resultTable := buildResultTable(results)
	h.Template.Execute(h.Writer, resultTable)
}
