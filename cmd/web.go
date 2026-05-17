package main

import (
	_ "embed"
	"os"

	"github.com/gofiber/fiber/v3"
)

//go:embed static/index.html
var staticFS []byte

//go:embed static/logo.png
var logoPNG []byte

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
		return c.Send(staticFS)
	})

	// Logo
	app.Get("/logo.png", func(c fiber.Ctx) error {
		c.Set("Content-Type", "image/png")
		return c.Send(logoPNG)
	})

	// 版本号
	app.Get("/api/version", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"version": Version})
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
		if err := SaveCurlFile("session.txt", req.Content); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Save failed"})
		}
		return c.JSON(fiber.Map{"ok": true})
	})

	// 获取最新消息
	app.Get("/api/poll", func(c fiber.Ctx) error {
		config, err := GetConfig("session.txt")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		msgID, msg, err := GetLatestMessage(config)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"ok": true, "message_id": msgID, "message": msg})
	})

	app.Listen(addr)
}

// StopWeb 关闭 Web 服务
func StopWeb() {
	if webApp != nil {
		webApp.Shutdown()
		webApp = nil
	}
}
