package main

import (
	_ "embed"
	_ "image/png"

	"fmt"
	"os"

	"fyne.io/systray"
)

//go:embed static/logo.png
var logo2PNG []byte

func StartTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	icoData, err := pngToICO(logo2PNG)
	if err != nil {
		fmt.Println("图标转换失败:", err)
		icoData = logo2PNG
	}
	systray.SetIcon(icoData)
	systray.SetTitle("豆包语音输入")
	systray.SetTooltip("豆包语音输入")

	mOpen := systray.AddMenuItem("打开设置页面", "打开浏览器进行配置")
	mClose := systray.AddMenuItem("关闭设置页面", "关闭配置页面并停止 Web 服务")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出程序")

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				// 启动 Web 服务
				addr := ":2828"
				if p := os.Getenv("DOUBAO_INPUT_PORT"); p != "" {
					addr = ":" + p
				}
				go StartWeb(addr)
				fmt.Printf("Web 界面: http://localhost%s\n", addr)
				openBrowser("http://localhost:2828")
			case <-mClose.ClickedCh:
				// 关闭 Web 服务
				StopWeb()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// clean up here
}
