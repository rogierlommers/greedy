package selfdiagnose

import (
	"os"
	"testing"
)

func TestCheckDirectory(t *testing.T) {
	testfile, err := os.Create("/tmp/selfdiagnose_test")
	if err != nil {
		t.Fatal("unable to create testfile:" + err.Error())
	}
	testfile.Close()

	Register(CheckDirectory{Path: "/tmp", CanAppend: true})
	Register(CheckDirectory{Path: "/tmp/missing", CanAppend: true})
	Register(CheckDirectory{Path: "/tmp/selfdiagnose_test", CanAppend: true})
	Run(LoggingReporter{})
	Run(HtmlReporter{os.Stdout})
}

func ExampleCheckDirectory() {
	check := CheckDirectory{Path: "/tmp", CanAppend: true}
	check.SetComment("something")
	Register(check)
}
