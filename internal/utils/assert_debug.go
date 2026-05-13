//go:build assertions

package utils

import "fmt"

// When the "assertions" flag is present, failed assertions panic.
//
// Otherwise, [Assert] is a no-op.
func Assert(cond bool) {
	if !cond {
		panic(fmt.Errorf("assertion failed!"))
	}
}
