package system

import (
	_ "embed"
	_ "image/png"

	"fmt"

	"Doubao-Susurro/assets"
	"Doubao-Susurro/internal/config"
	"Doubao-Susurro/internal/system/startup"
	"Doubao-Susurro/internal/tool"

	"github.com/energye/systray"
)

func StartTray() {
	systray.Run(onReady, onExit)
}

func openSetting() {
	tool.OpenBrowser()
}

var mAutoType *systray.MenuItem
var mStartup *systray.MenuItem

func taggleAutoType() {
	cfg := config.GetConfig()
	cfg.AutoType = !cfg.AutoType
	config.SaveConfig(cfg)
	if config.GetConfig().AutoType {
		mAutoType.Check()
	} else {
		mAutoType.Uncheck()
	}
}

func taggleStartup() {
	cfg := config.GetConfig()
	newStartupState := !cfg.Startup

	if err := startup.UpdateStartup(newStartupState); err != nil {
		fmt.Println("更新开机自启状态失败:", err)
		return
	}

	// 只有操作成功后才更新配置和UI
	cfg.Startup = newStartupState
	config.SaveConfig(cfg)

	// 更新UI状态
	if newStartupState {
		mStartup.Check()
	} else {
		mStartup.Uncheck()
	}
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

	mAutoType = systray.AddMenuItemCheckbox("自动输入", "启用自动输入功能", config.GetConfig().AutoType)
	mAutoType.Click(taggleAutoType)
	mStartup = systray.AddMenuItemCheckbox("开机自启", "启用开机自启功能", config.GetConfig().Startup)
	mStartup.Click(taggleStartup)
	mOpen := systray.AddMenuItem("设置", "打开设置页面")
	mOpen.Click(openSetting)
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出程序")
	mQuit.Click(systray.Quit)
}

func onExit() {
	// clean up here
}
