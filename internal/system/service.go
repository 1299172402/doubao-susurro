package system

import (
	"Doubao-input/internal/core"
	"flag"
	"fmt"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	core.StartClipboardWriter()
}

func (p *program) Stop(s service.Service) error {
	return nil
}

// RunService 解析命令行参数并管理 Windows 服务的安装、卸载、启动和运行。
func RunService() {
	var install, uninstall, start, stop bool
	flag.BoolVar(&install, "install", false, "安装服务")
	flag.BoolVar(&uninstall, "uninstall", false, "卸载服务")
	flag.BoolVar(&start, "start", false, "启动服务")
	flag.BoolVar(&stop, "stop", false, "停止服务")
	flag.Parse()

	srvConfig := &service.Config{
		Name:        "DoubaoInput",
		DisplayName: "豆包语音输入",
		Description: "豆包语音输入",
	}
	prg := &program{}
	s, err := service.New(prg, srvConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	if install || uninstall || start || stop {
		switch {
		case install:
			if err := s.Install(); err != nil {
				fmt.Println("安装服务失败:", err.Error())
			}
		case uninstall:
			if err := s.Uninstall(); err != nil {
				fmt.Println("卸载服务失败:", err.Error())
			}
		case start:
			if err := s.Start(); err != nil {
				fmt.Println("运行服务失败:", err.Error())
			}
		case stop:
			if err := s.Stop(); err != nil {
				fmt.Println("停止服务失败:", err.Error())
			}
		}
		return
	}

	if err := s.Run(); err != nil {
		fmt.Println(err)
	}
}
