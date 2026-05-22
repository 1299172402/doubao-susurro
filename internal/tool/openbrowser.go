package tool

import (
	"Doubao-Susurro/internal/config"
	"fmt"
	"os/exec"
	"runtime"
)

// openBrowser 在不同平台上打开默认浏览器访问指定 URL
func OpenBrowser() {
	url := fmt.Sprintf("http://localhost:%s", config.GetConfig().Port)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start()
}
