//go:build windows

package update

import "syscall"

func detachedProc() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: 0x00000008} // DETACHED_PROCESS
}
