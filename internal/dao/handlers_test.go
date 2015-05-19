package dao

import (
	"fmt"
)

func testIsAboutBlank() {
	aboutTests := []string{"YWJvdXQ6Ymxhbms=", "http://www.lommers.org"}

	for _, value := range aboutTests {
		check := isAboutBlank(value)
	}

	fmt.Println(check)
	// Output: golly
}
