// +build !aix,!darwin,!dragonfly,!freebsd,!linux,!nacl,!nacljs,!netbsd,!openbsd,!solaris,!windows

// For systems without syscall.Errno.
// Build targets must be inverse of errors_errno.go

package piondtlsdynamo

import (
	"os"
)

func isOpErrorTemporary(err *os.SyscallError) bool {
	return false
}
