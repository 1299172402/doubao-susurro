package main

import (
	"fmt"
)

// 通过 ldflags 在构建时注入，例如:
// go build -ldflags="-X main.Version=v1.0.0" -o doubao-input.exe .
var Version = "dev"

func main() {
	fmt.Printf("Doubao Input %s\n", Version)

	// 启动消息监听（后台运行，不阻塞）
	go StartClipboardWriter()

	// 启动系统托盘（阻塞 main）
	StartTray()
}
