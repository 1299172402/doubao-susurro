package startup

import (
	"fmt"
	"os"
	"path/filepath"

	"Doubao-input/assets"
)

// Windows 实现

func installWindowsStartup() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取执行文件路径失败: %v", err)
	}

	appData := os.Getenv("APPDATA")
	if appData == "" {
		return fmt.Errorf("无法获取 APPDATA 环境变量")
	}

	startupDir := filepath.Join(appData, "Microsoft\\Windows\\Start Menu\\Programs\\Startup")
	vbsContent := fmt.Sprintf(assets.StartupWindowsVBS, exePath)
	vbsPath := filepath.Join(startupDir, "doubao-input-startup.vbs")

	if err := os.WriteFile(vbsPath, []byte(vbsContent), 0644); err != nil {
		return fmt.Errorf("写入 VBS 文件失败: %v", err)
	}

	return nil
}

func uninstallWindowsStartup() error {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return fmt.Errorf("无法获取 APPDATA 环境变量")
	}

	startupDir := filepath.Join(appData, "Microsoft\\Windows\\Start Menu\\Programs\\Startup")
	vbsPath := filepath.Join(startupDir, "doubao-input-startup.vbs")

	if err := os.Remove(vbsPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除 VBS 文件失败: %v", err)
	}

	return nil
}
