package main

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
)

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
