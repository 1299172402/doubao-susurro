package main

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
)

// 通过 ldflags 在构建时注入，例如:
// go build -ldflags="-X main.Version=v1.0.0" -o doubao-input.exe .
var Version = "dev"

func main() {
	config, err := GetConfig("session.txt")
	if err != nil {
		fmt.Printf("配置加载失败: %v\n", err)
		return
	}

	fmt.Printf("Doubao Input %s\n", Version)
	fmt.Println("开始监听豆包消息...")
	fmt.Println("按 Ctrl+C 停止")

	for {
		msg, err := PollNewMessage(config)
		if err != nil {
			fmt.Printf("轮询错误: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if msg != "" {
			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("[%s] %s\n", now, msg)

			if err := clipboard.WriteAll(msg); err != nil {
				fmt.Printf("复制到剪贴板失败: %v\n", err)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
