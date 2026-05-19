package core

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
)

// 记录上次消息 ID，避免重复处理
var lastMessageID string

func StartPolling() {
	for {
		msgID, msg, err := DeliverMessage()
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
			robotgo.Type(msg)
		}

		time.Sleep(1 * time.Second)
	}
}
