//go:build windows

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// isAdmin 检查当前进程是否以管理员权限运行
func isAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}
	return member
}

// elevateAndRestart 以管理员权限重新启动当前进程。
// 通过 ShellExecuteW 的 "runas" 动词触发 Windows UAC 提权弹窗。
func elevateAndRestart() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.Abs(exe)
	if err != nil {
		return err
	}

	shell32 := windows.NewLazyDLL("shell32.dll")
	proc := shell32.NewProc("ShellExecuteW")

	verbPtr, _ := windows.UTF16PtrFromString("runas")
	exePtr, _ := windows.UTF16PtrFromString(exe)

	args := ""
	if len(os.Args) > 1 {
		args = strings.Join(os.Args[1:], " ")
	}
	argsPtr, _ := windows.UTF16PtrFromString(args)

	// SW_SHOWNORMAL = 1
	ret, _, callErr := proc.Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(exePtr)),
		uintptr(unsafe.Pointer(argsPtr)),
		0,
		1,
	)
	if ret <= 32 {
		return fmt.Errorf("ShellExecuteW 失败: %v", callErr)
	}
	return nil
}

// ensureAdmin 检查是否以管理员权限运行，如果不是则自动提权并退出当前进程。
// 返回 true 表示已触发提权（调用方应立即退出），false 表示已是管理员。
func ensureAdmin() bool {
	if isAdmin() {
		return false
	}

	fmt.Println("需要管理员权限，正在请求提权...")
	if err := elevateAndRestart(); err != nil {
		fmt.Println("自动提权失败:", err)
		return false
	}

	os.Exit(0)
	return true // unreachable
}
