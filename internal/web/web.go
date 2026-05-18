package web

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"

	"Doubao-input/assets"
	"Doubao-input/info"
	"Doubao-input/internal/core"
	"Doubao-input/internal/tool"
)

var webApp *fiber.App

// StartWeb 启动 Web 服务
func StartWeb(addr string) {
	if webApp != nil {
		return // 已经在运行
	}

	app := fiber.New()
	webApp = app

	// 首页
	app.Get("/", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Send(assets.IndexPage)
	})

	// Logo
	app.Get("/logo.png", func(c fiber.Ctx) error {
		c.Set("Content-Type", "image/png")
		return c.Send(assets.LogoPNG)
	})

	// 版本号
	app.Get("/api/version", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"version": info.Version})
	})

	// 获取 session
	app.Get("/api/session", func(c fiber.Ctx) error {
		data, err := os.ReadFile("session.txt")
		if err != nil {
			return c.JSON(fiber.Map{"content": ""})
		}
		return c.JSON(fiber.Map{"content": string(data)})
	})

	// 保存 session
	app.Post("/api/session/save", func(c fiber.Ctx) error {
		var req struct {
			Content string `json:"content"`
		}
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
		}
		if err := tool.WriteCurlFile("session.txt", req.Content); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Save failed"})
		}
		return c.JSON(fiber.Map{"ok": true})
	})

	// 获取最新消息
	app.Get("/api/poll", func(c fiber.Ctx) error {
		msgID, msg, err := core.DeliverMessage()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"ok": true, "message_id": msgID, "message": msg})
	})

	// 只监听 127.0.0.1，禁止远程访问
	bindAddr := "127.0.0.1" + addr
	app.Listen(bindAddr)
}

// StopWeb 关闭 Web 服务
func StopWeb() {
	if webApp != nil {
		webApp.Shutdown()
		webApp = nil
	}
}

func Launch() {
	// 启动 Web 服务
	addr := ":2828"
	if p := os.Getenv("DOUBAO_INPUT_PORT"); p != "" {
		addr = ":" + p
	}
	go StartWeb(addr)
	// 自动打开浏览器
	fmt.Printf("Web 界面: http://localhost%s\n", addr)
	tool.OpenBrowser(fmt.Sprintf("http://localhost%s", addr))
}
