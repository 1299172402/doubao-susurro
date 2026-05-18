package lock

import (
	"fmt"
	"os"
	"path/filepath"
)

// TryLock 尝试获取单实例锁，返回 unlock 函数和 error
func TryLock(name string) (unlock func(), err error) {
	lockPath := filepath.Join(os.TempDir(), name+".lock")

	// 尝试创建并锁定文件
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("无法创建锁文件: %w", err)
	}

	// 尝试加排他锁，非阻塞
	if err := lockFile(f); err != nil {
		f.Close()
		return nil, fmt.Errorf("程序已在运行中")
	}

	// 写入当前 PID
	f.Truncate(0)
	f.Seek(0, 0)
	fmt.Fprintf(f, "%d", os.Getpid())
	f.Sync()

	unlock = func() {
		unlockFile(f)
		f.Close()
		os.Remove(lockPath)
	}

	return unlock, nil
}
