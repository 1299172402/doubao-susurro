package system

import (
	_ "embed"
	_ "image/png"

	"fmt"

	"Doubao-input/assets"
	"Doubao-input/internal/tool"
	"Doubao-input/internal/web"

	"github.com/energye/systray"
)

func StartTray() {
	systray.Run(onReady, onExit)
}

func openSetting() {
	web.Launch()
}

func onReady() {
	icoData, err := tool.PngToIco(assets.LogoPNG)
	if err != nil {
		fmt.Println("图标转换失败:", err)
		icoData = assets.LogoPNG
	}
	systray.SetIcon(icoData)
	systray.SetTitle("豆包语音输入")
	systray.SetTooltip("豆包语音输入")

	mOpen := systray.AddMenuItem("打开设置页面", "打开浏览器进行配置")
	mOpen.Click(openSetting)
	mClose := systray.AddMenuItem("关闭设置页面", "关闭配置页面并停止 Web 服务")
	mClose.Click(web.StopWeb)

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出程序")
	mQuit.Click(systray.Quit)
}

func onExit() {
	// clean up here
}
