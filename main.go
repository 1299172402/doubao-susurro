package main

import (
	"fmt"
	"os"
	"time"

	"github.com/atotto/clipboard"
)

// 通过 ldflags 在构建时注入，例如:
// go build -ldflags="-X main.Version=v1.0.0" -o doubao-input.exe .
var Version = "dev"

// 记录上次消息 ID，避免重复处理
var lastMessageID string

func StartClipboardWriter() {
	for {
		config, err := GetConfig("session.txt")
		if err != nil {
			fmt.Println("配置加载失败:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		msgID, msg, err := GetLatestMessage(config)
		if err != nil {
			fmt.Println("轮询错误:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if msgID != lastMessageID {
			lastMessageID = msgID

			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("[%s] %s\n", now, msg)
			clipboard.WriteAll(msg)
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	fmt.Printf("Doubao Input %s\n", Version)

	// 启动 Web 服务
	addr := ":2828"
	if p := os.Getenv("DOUBAO_INPUT_PORT"); p != "" {
		addr = ":" + p
	}
	go StartWeb(addr)
	fmt.Printf("Web 界面: http://localhost%s\n", addr)

	// 启动消息监听和剪贴板写入
	StartClipboardWriter()
}
