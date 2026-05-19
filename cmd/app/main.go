package main

import (
	"Doubao-input/internal/config"
	"Doubao-input/internal/core"
	"Doubao-input/internal/system"
	"Doubao-input/internal/system/lock"
	"Doubao-input/internal/tool"
	"Doubao-input/internal/web"
	"flag"
	"fmt"
	"os"
)

func main() {
	silent := flag.Bool("silent", false, "静默模式，不打开浏览器")
	flag.Parse()

	// 确保单实例运行
	unlock, err := lock.TryLock("doubao-input")
	if err != nil {
		fmt.Println(err)
		tool.OpenBrowser()
		os.Exit(1)
	}
	defer unlock()

	// 加载配置
	config.InitConfig()

	// 启动消息监听（后台运行，不阻塞）
	go core.StartPolling()

	// 启动网页设置服务器（后台运行，不阻塞）
	go web.StartWeb()

	// 非静默模式
	if !*silent {
		// 启动设置页面服务器和自动打开浏览器
		tool.OpenBrowser()
	}

	// 启动系统托盘（阻塞 main）
	system.StartTray()
}
