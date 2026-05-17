package main

import (
	"fmt"
	"runtime"

	"github.com/kardianos/service"
)

type Program struct{}

func (p *Program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *Program) run() {
	fmt.Println("服务运行中...")
	// 启动消息监听
	go StartClipboardWriter()
	// 启动系统托盘
	go StartTray()
}
func (p *Program) Stop(s service.Service) error {
	fmt.Println("服务停止")
	// 停止 Web 服务
	StopWeb()
	return nil
}

func initService() service.Service {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg := &service.Config{
		Name:        "doubao-input",
		DisplayName: "Doubao Input",
		Description: "Doubao Input Service",
	}
	prg := &Program{}
	s, _ := service.New(prg, cfg)
	return s
}

func installService() {
	s := initService()
	err := s.Install()
	if err != nil {
		fmt.Println("安装服务失败:", err)
	}
}

func uninstallService() {
	s := initService()
	err := s.Uninstall()
	if err != nil {
		fmt.Println("卸载服务失败:", err)
	}
}

func isServiceInstalled() bool {
	s := initService()

	_, err := s.Status()
	if err != nil {
		return false
	}
	return true
}
