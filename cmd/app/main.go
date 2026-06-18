package main

import (
	"Doubao-Susurro/internal/config"
	"Doubao-Susurro/internal/core"
	"Doubao-Susurro/internal/system"
	"Doubao-Susurro/internal/system/lock"
	"Doubao-Susurro/internal/tool"
	"Doubao-Susurro/internal/web"
	"flag"
	"fmt"
	"os"
)

func main() {
	silent := flag.Bool("silent", false, "静默模式，不打开浏览器")
	configPath := flag.String("config", "", "配置文件路径，默认可执行文件所在目录下的 Doubao-Susurro-config.yml")
	flag.Parse()

	// 加载配置
	config.InitConfig(*configPath)

	// 确保单实例运行
	unlock, err := lock.TryLock("Doubao-Susurro")
	if err != nil {
		fmt.Println(err)
		tool.OpenBrowser()
		os.Exit(1)
	}
	defer unlock()

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
