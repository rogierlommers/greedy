package selfdiagnose

import (
	"testing"
	"time"
)

type long struct {
	BasicTask
}

func (l long) Run(ctx *Context, result *Result) {
	time.Sleep(10 * time.Second)
	result.Passed = true
	result.Reason = "waited out"
}
func TestThatLongTasksAreTimedout(t *testing.T) {
	l := new(long)
	l.SetTimeout(5 * time.Second)
	Register(l)
	Run(LoggingReporter{})
}
