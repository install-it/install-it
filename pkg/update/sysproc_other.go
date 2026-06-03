//go:build !windows

package update

import "syscall"

func detachedProc() *syscall.SysProcAttr {
	return nil
}
