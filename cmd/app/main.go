package main

import (
	"flag"
	"fmt"
	"os"

	"Doubao-input/info"
	"Doubao-input/internal/core"
	"Doubao-input/internal/system"
	"Doubao-input/internal/system/lock"
	"Doubao-input/internal/web"
)

func main() {
	silent := flag.Bool("silent", false, "静默模式，不打开浏览器")
	flag.Parse()

	// 防止重复运行
	unlock, err := lock.TryLock("doubao-input")
	if err != nil {
		fmt.Println(err)
		web.Launch()
		os.Exit(1)
	}
	defer unlock()

	// 交互式模式（双击启动）
	fmt.Printf("Doubao Input\n")
	fmt.Printf("Version: %s\n", info.Version)

	// 启动消息监听（后台运行，不阻塞）
	go core.StartClipboardWriter()

	// 非静默模式
	if !*silent {
		web.Launch()
	}

	// 启动系统托盘（阻塞 main）
	system.StartTray()
}
