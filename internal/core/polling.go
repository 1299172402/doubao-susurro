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
	for {
		msgID, msg, err := GetLatestMessage()
		if err != nil {
			fmt.Println("轮询错误:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if msgID != lastMessageID {
			lastMessageID = msgID

			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("[%s] %s\n", now, msg)

			// 复制到剪贴板
			clipboard.WriteAll(msg)

			// 自动输入
			if config.GetConfig().AutoType {
				robotgo.Type(msg)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
