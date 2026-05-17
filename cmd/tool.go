package main

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
	"image"
	_ "image/png"
	"os/exec"
	"runtime"
)

// pngToICO 将 PNG 字节转换为 ICO 格式
func pngToICO(pngData []byte) ([]byte, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(pngData))
	if err != nil {
		return nil, fmt.Errorf("解析 PNG 尺寸失败: %w", err)
	}

	width := uint8(cfg.Width)
	height := uint8(cfg.Height)
	if cfg.Width > 255 {
		width = 0
	}
	if cfg.Height > 255 {
		height = 0
	}

	dataSize := uint32(len(pngData))
	offset := uint32(6 + 16)

	var buf bytes.Buffer
	// ICONDIR
	buf.Write([]byte{0, 0})
	binary.Write(&buf, binary.LittleEndian, uint16(1))
	binary.Write(&buf, binary.LittleEndian, uint16(1))
	// ICONDIRENTRY
	buf.WriteByte(width)
	buf.WriteByte(height)
	buf.WriteByte(0)
	buf.WriteByte(0)
	binary.Write(&buf, binary.LittleEndian, uint16(1))
	binary.Write(&buf, binary.LittleEndian, uint16(32))
	binary.Write(&buf, binary.LittleEndian, dataSize)
	binary.Write(&buf, binary.LittleEndian, offset)
	// PNG data
	buf.Write(pngData)

	return buf.Bytes(), nil
}

// openBrowser 在不同平台上打开默认浏览器访问指定 URL
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start()
}
