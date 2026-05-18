package main

import (
	"fmt"

	"Doubao-input/info"
	"Doubao-input/internal/system"
)

func main() {
	fmt.Printf("Doubao Input\n")
	fmt.Printf("Version: %s\n", info.Version)

	// 启动服务
	system.RunService()
}
