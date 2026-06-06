package core

import (
	"Doubao-input/internal/config"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
)

// 记录上次消息 ID，避免重复处理
var lastMessageID string

func StartPolling() {
	// 启动时跳过最近一条消息，避免一启动就自动输入旧消息
	firstRun := true

	for {
		msgID, msg, err := GetLatestMessage()
		if err != nil {
			fmt.Println("轮询错误:", err)
			time.Sleep(time.Duration(config.GetConfig().IntervalTime) * time.Millisecond)
			continue
		}

		if msgID != lastMessageID {
			lastMessageID = msgID

			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("[%s] %s\n", now, msg)

			// 复制到剪贴板
			clipboard.WriteAll(msg)

			// 自动输入
			if config.GetConfig().AutoType && !firstRun {
				robotgo.KeyTap("v", "ctrl")
			}

			firstRun = false

		}

		time.Sleep(time.Duration(config.GetConfig().IntervalTime) * time.Millisecond)
	}
}
