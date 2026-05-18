//go:build windows

package lock

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx = kernel32.NewProc("LockFileEx")
)

const (
	LOCKFILE_EXCLUSIVE_LOCK   = 0x00000002
	LOCKFILE_FAIL_IMMEDIATELY = 0x00000001
)

type overlapped struct {
	Internal     uintptr
	InternalHigh uintptr
	Offset       uint32
	OffsetHigh   uint32
	HEvent       uintptr
}

func lockFile(f *os.File) error {
	h := syscall.Handle(f.Fd())
	var o overlapped
	r, _, err := procLockFileEx.Call(
		uintptr(h),
		uintptr(LOCKFILE_EXCLUSIVE_LOCK|LOCKFILE_FAIL_IMMEDIATELY),
		0,
		1,
		0,
		uintptr(unsafe.Pointer(&o)),
	)
	if r == 0 {
		return err
	}
	return nil
}

func unlockFile(f *os.File) error {
	h := syscall.Handle(f.Fd())
	var o overlapped
	r, _, err := procLockFileEx.Call(
		uintptr(h),
		0,
		0,
		1,
		0,
		uintptr(unsafe.Pointer(&o)),
	)
	if r == 0 {
		return err
	}
	return nil
}
