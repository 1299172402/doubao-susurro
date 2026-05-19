package startup

import (
	"Doubao-input/assets"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// macOS 实现

func installMacOSStartup() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取执行文件路径失败: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %v", err)
	}

	launchAgentsDir := filepath.Join(homeDir, "Library/LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return fmt.Errorf("创建 LaunchAgents 目录失败: %v", err)
	}

	plistContent := fmt.Sprintf(assets.StartupMacOSPlist, exePath)
	plistPath := filepath.Join(launchAgentsDir, "com.doubao.input.plist")

	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		return fmt.Errorf("写入 plist 文件失败: %v", err)
	}

	// 加载 launchd 任务
	cmd := exec.Command("launchctl", "load", "-w", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("加载启动任务失败: %v, 输出: %s", err, output)
	}

	return nil
}

func uninstallMacOSStartup() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %v", err)
	}

	plistPath := filepath.Join(homeDir, "Library/LaunchAgents", "com.doubao.input.plist")

	// 卸载任务（忽略错误，可能未加载）
	exec.Command("launchctl", "unload", "-w", plistPath).Run()

	// 删除文件
	if err := os.Remove(plistPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除 plist 文件失败: %v", err)
	}

	return nil
}
