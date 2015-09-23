package selfdiagnose

import (
	"fmt"
	"os"
)

// CheckDirectory reports whether a directory (path) is readable and/or appendable
type CheckDirectory struct {
	BasicTask
	Path      string
	CanAppend bool
}

func (c CheckDirectory) Run(ctx *Context, result *Result) {
	result.Passed = false
	info, err := os.Stat(c.Path)
	if os.IsNotExist(err) {
		result.Reason = fmt.Sprintf("Directory [%s] does not exist.", c.Path)
		return
	}
	if err != nil {
		result.Reason = err.Error()
		return
	}
	if !info.IsDir() {
		result.Reason = fmt.Sprintf("Directory [%s] exists but is not a directory. Permission:[%v]", c.Path, info.Mode())
		return
	}
	if info.Mode()&os.ModeAppend != 0 {
		result.Reason = fmt.Sprintf("Directory [%s] exists but files cannot be appended. Permission:[%v]", c.Path, info.Mode())
		return
	}
	result.Passed = true
	result.Reason = fmt.Sprintf("Directory [%s] exists and is appendable", c.Path)
}
