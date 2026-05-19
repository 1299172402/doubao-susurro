package startup

import (
	"Doubao-input/assets"
	"fmt"
	"os"
	"path/filepath"
)

// Linux 实现
func installLinuxStartup() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取执行文件路径失败: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %v", err)
	}

	autostartDir := filepath.Join(homeDir, ".config", "autostart")
	if err := os.MkdirAll(autostartDir, 0755); err != nil {
		return fmt.Errorf("创建 autostart 目录失败: %v", err)
	}

	desktopContent := fmt.Sprintf(assets.StartupLinuxDesktop, exePath)
	desktopPath := filepath.Join(autostartDir, "doubao-input.desktop")

	if err := os.WriteFile(desktopPath, []byte(desktopContent), 0644); err != nil {
		return fmt.Errorf("写入 desktop 文件失败: %v", err)
	}

	return nil
}

func uninstallLinuxStartup() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %v", err)
	}

	desktopPath := filepath.Join(homeDir, ".config", "autostart", "doubao-input.desktop")

	if err := os.Remove(desktopPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除 desktop 文件失败: %v", err)
	}

	return nil
}
