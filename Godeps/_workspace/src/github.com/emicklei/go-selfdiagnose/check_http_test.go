package selfdiagnose

import (
	"net/http"
	"testing"
)

func TestCheckHttp(t *testing.T) {
	get, err := http.NewRequest("GET", "http://ernestmicklei.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	check := CheckHttp{Request: get}
	check.SetComment("blog access")
	Register(check)
	Run(LoggingReporter{})
}

func ExampleCheckHttp() {
	get, _ := http.NewRequest("GET", "http://ernestmicklei.com", nil)
	check := CheckHttp{Request: get}
	check.SetComment("blog access")
	Register(check)
}
