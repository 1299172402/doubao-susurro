package main

import (
	"flag"
	"fmt"
	"os"
)

// 通过 ldflags 在构建时注入，例如:
// go build -ldflags="-X main.Version=v1.0.0" -o doubao-input.exe .
var Version = "dev"

func main() {
	silent := flag.Bool("silent", false, "静默模式，不打开浏览器")
	flag.Parse()

	// 交互式模式（双击启动）
	fmt.Printf("Doubao Input\n")
	fmt.Printf("Version: %s\n", Version)

	// 启动消息监听（后台运行，不阻塞）
	go StartClipboardWriter()

	// 非静默模式
	if !*silent {
		// 启动 Web 服务
		addr := ":2828"
		if p := os.Getenv("DOUBAO_INPUT_PORT"); p != "" {
			addr = ":" + p
		}
		go StartWeb(addr)
		// 自动打开浏览器
		go openBrowser("http://localhost:2828")
	}

	// 启动系统托盘（阻塞 main）
	StartTray()
}
