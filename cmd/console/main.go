package main

import (
	"fmt"

	"Doubao-input/info"
	"Doubao-input/internal/core"
)

func main() {
	fmt.Printf("Doubao Input\n")
	fmt.Printf("Version: %s\n", info.Version)

	core.StartClipboardWriter()
}
