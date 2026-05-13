//go:build !assertions

package utils

// When the "assertions" flag is present, failed assertions panic.
//
// When disabled, [Assert] is a no-op.
func Assert(cond bool) {}
