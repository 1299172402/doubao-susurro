package web

import (
	_ "embed"

	"github.com/gofiber/fiber/v3"

	"Doubao-input/assets"
	"Doubao-input/info"
	"Doubao-input/internal/config"
	"Doubao-input/internal/core"
	"Doubao-input/internal/system/startup"
)

var webApp *fiber.App

// StartWeb 启动 Web 服务
func StartWeb() {
	if webApp != nil {
		return // 已经在运行
	}

	// 启动 Web 服务
	port := config.GetConfig().Port
	addr := ":" + port

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

	// 捐赠图片
	app.Get("/donate/wechat.png", func(c fiber.Ctx) error {
		c.Set("Content-Type", "image/png")
		return c.Send(assets.DonateWechat)
	})
	
	app.Get("/donate/alipay.jpg", func(c fiber.Ctx) error {
		c.Set("Content-Type", "image/jpeg")
		return c.Send(assets.DonateAlipay)
	})

	// 版本号
	app.Get("/api/version", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"version": info.Version})
	})

	// 获取 session
	app.Get("/api/session", func(c fiber.Ctx) error {
		data := config.GetConfig().Session
		return c.JSON(fiber.Map{"content": string(data)})
	})

	// 获取所有配置
	app.Get("/api/config", func(c fiber.Ctx) error {
		cfg := config.GetConfig()
		return c.JSON(fiber.Map{
			"ok":                 true,
			"auto_type":          cfg.AutoType,
			"conversation_limit": cfg.ConversationLimit,
			"interval_time":      cfg.IntervalTime,
			"startup":            cfg.Startup,
			"port":               cfg.Port,
		})
	})

	// 保存配置
	app.Post("/api/config/save", func(c fiber.Ctx) error {
		var req struct {
			AutoType          *bool `json:"auto_type"`
			ConversationLimit *int  `json:"conversation_limit"`
			IntervalTime      *int  `json:"interval_time"`
			Startup           *bool `json:"startup"`
		}
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
		}

		cfg := config.GetConfig()
		if req.AutoType != nil {
			cfg.AutoType = *req.AutoType
		}
		if req.ConversationLimit != nil {
			cfg.ConversationLimit = *req.ConversationLimit
		}
		if req.IntervalTime != nil {
			cfg.IntervalTime = *req.IntervalTime
		}
		if req.Startup != nil {
			newStartupState := *req.Startup

			if err := startup.UpdateStartup(newStartupState); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "更新开机自启状态失败: " + err.Error()})
			}

			cfg.Startup = newStartupState
		}

		if err := config.SaveConfig(cfg); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Save failed"})
		}
		return c.JSON(fiber.Map{"ok": true})
	})

	// 保存 session
	app.Post("/api/session/save", func(c fiber.Ctx) error {
		var req struct {
			Content string `json:"content"`
		}
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
		}

		cfg := config.GetConfig()
		cfg.Session = req.Content
		if err := config.SaveConfig(cfg); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Save failed"})
		}

		return c.JSON(fiber.Map{"ok": true})
	})

	// 获取最新消息
	app.Get("/api/poll", func(c fiber.Ctx) error {
		msgID, msg, err := core.GetLatestMessage()
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
