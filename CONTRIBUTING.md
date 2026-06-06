
# For Developers

## 环境

- Go 1.25+
- GCC（`robotgo` 依赖 CGO）

## 运行

```bash
go run ./cmd/app
```

## 项目结构

```
├── cmd/
│   └── app/
│       └── main.go                   # 主程序入口
├── internal/
│   ├── config/
│   │   └── config.go                 # 配置管理（Viper，YAML 格式）
│   ├── core/
│   │   ├── curl_parser.go            # cURL 命令解析
│   │   ├── listener.go               # 调用豆包接口获取最新消息
│   │   └── polling.go                # 消息轮询主循环，自动复制到剪贴板
│   ├── system/
│   │   ├── tray.go                   # 系统托盘菜单管理
│   │   ├── lock/
│   │   │   ├── lock.go               # 进程锁接口（跨平台）
│   │   │   ├── lock_windows.go       # Windows 文件锁实现（LockFileEx）
│   │   │   └── lock_unix.go          # Unix 文件锁实现（flock）
│   │   └── startup/
│   │       ├── startup.go            # 开机自启接口（跨平台）
│   │       ├── windows.go            # Windows 开机自启（VBS 脚本）
│   │       ├── linux.go              # Linux 开机自启（.desktop 文件）
│   │       └── macos.go              # macOS 开机自启（LaunchAgent plist）
│   ├── tool/
│   │   ├── openbrowser.go            # 跨平台打开默认浏览器
│   │   └── pngtoico.go               # PNG 转 ICO 格式
│   └── web/
│       └── web.go                    # Web 服务（Fiber 框架）
├── assets/
│   ├── asset.go                      # go:embed 资源声明
│   └── static/
│       ├── index.html                # Web 设置界面
│       ├── logo.png                  # 应用图标
│       └── startup/
│           ├── windows.vbs           # Windows 开机自启脚本模板
│           ├── linux.desktop         # Linux 开机自启桌面文件模板
│           └── macos.plist           # macOS LaunchAgent 配置模板
├── info/
│   └── version.go                    # 版本号定义（构建时注入）
├── Doubao-Susurro-config.yml           # 运行时生成的配置文件
└── go.mod
```

## 前端说明

Web 设置界面（`assets/static/index.html`）为单页应用，主要功能：

| 区域 | 说明 |
|------|------|
| Session 输入框 | 粘贴从豆包复制的 cURL 命令 |
| 保存配置按钮 | 保存 Session 内容 |
| 获取消息按钮 | 手动拉取一次最新消息 |
| 自动输入开关 | 切换后自动保存 |
| 开机自启开关 | 切换后自动保存 |
| 获取数量 | 数字输入，修改后自动保存 |
| 对话 ID | 文本输入，留空则自动从 curl 中提取；修改后自动保存 |
| 请求间隔 | 数字输入，修改后自动保存 |

所有设置项修改后均通过 `onchange` 事件**自动保存**到后端，并显示 Toast 提示。

## 模块说明

### `cmd/app/main.go`

主程序入口。启动流程：

1. 解析命令行参数（`-silent`）
2. 加载配置文件 `Doubao-Susurro-config.yml`（不存在则自动创建）
3. 获取单实例锁（防止重复启动）
4. 后台启动消息轮询（`core.StartPolling`）
5. 后台启动 Web 服务（`web.StartWeb`）
6. 非静默模式下自动打开浏览器
7. 启动系统托盘（阻塞主线程）

### `internal/config/config.go`

基于 [Viper](https://github.com/spf13/viper) 的配置管理。配置文件为 YAML 格式，环境变量前缀为 `DOUBAO_SUSURRO`。

```go
type Config struct {
    Port              string `mapstructure:"port"`               // Web 服务端口，默认 "2828"
    AutoType          bool   `mapstructure:"auto_type"`          // 自动输入模式，默认 true
    Startup           bool   `mapstructure:"startup"`            // 开机自启，默认 false
    Session           string `mapstructure:"session"`            // 从豆包复制的 cURL 命令
    ConversationLimit int    `mapstructure:"conversation_limit"` // 单次获取对话数量，默认 5
    IntervalTime      int    `mapstructure:"interval_time"`      // 轮询间隔（毫秒），默认 1000
    ConversationID    string `mapstructure:"conversation_id"`    // 对话 ID，留空自动从 curl 中提取
}
```

- `InitConfig()` — 加载配置，不存在则创建默认配置文件
- `GetConfig()` — 获取全局配置实例
- `SaveConfig(cfg)` — 保存配置到文件

### `internal/core/curl_parser.go`

解析从豆包复制的 cURL 命令，提取 URL、请求头、Cookie、请求体等信息。

- `parseCurl(curlStr)` — 解析 cURL 字符串，返回 `curlConfig` 结构
- `getConfig()` — 从全局配置读取 session 并解析，自动修正请求参数：
  - `direction` → `0`（从旧到新）
  - `anchor_index` → `0`（从最新开始）
  - `limit` → 使用配置中的 `conversation_limit`
  - `conversation_id` → 如果配置了 `conversation_id` 则覆盖 body 和 `referer` 头中的对话 ID

### `internal/core/listener.go`

调用豆包接口获取最新一条用户消息。

- `GetLatestMessage()` — 发送 HTTP POST 请求，解析响应，返回 `(messageID, content, error)`，仅返回 `user_type == 1` 的用户消息

### `internal/core/polling.go`

消息轮询主循环，根据 `interval_time` 配置定时执行。

- 检测到新消息时自动写入系统剪贴板（`atotto/clipboard`）
- 如果 `auto_type` 为 `true`，通过 `robotgo` 自动将文本键入当前焦点窗口
- 通过 `lastMessageID` 去重，避免重复处理

### `internal/system/tray.go`

系统托盘菜单管理，使用 [energye/systray](https://github.com/energye/systray)。

- 托盘菜单项：自动输入（开关）、开机自启（开关）、设置（打开浏览器）、退出
- 托盘图标从嵌入的 PNG 资源实时转换为 ICO 格式

### `internal/system/startup/`

跨平台开机自启实现：

- **Windows**：在 `%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup` 创建 VBS 脚本
- **Linux**：在 `~/.config/autostart/` 创建 `.desktop` 文件
- **macOS**：在 `~/Library/LaunchAgents/` 创建 `.plist` 文件

### `internal/system/lock/`

单实例进程锁，防止程序重复启动：

- **Windows**：使用 `LockFileEx` API
- **Unix**：使用 `flock` 系统调用

### `internal/web/web.go`

基于 [Fiber v3](https://gofiber.io/) 的 Web 服务，仅监听 `127.0.0.1`。

API 端点：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/` | 设置页面（HTML） |
| GET | `/logo.png` | 应用图标 |
| GET | `/api/version` | 获取版本号 |
| GET | `/api/session` | 获取当前 session |
| POST | `/api/session/save` | 保存 session |
| GET | `/api/config` | 获取所有配置 |
| POST | `/api/config/save` | 保存配置 |
| GET | `/api/poll` | 手动获取最新消息 |
| GET | `/donate/wechat.png` | 微信收款码 |
| GET | `/donate/alipay.jpg` | 支付宝收款码 |

### `assets/asset.go`

使用 `go:embed` 嵌入静态资源：

- `IndexPage` — Web 设置页面 HTML
- `LogoPNG` — 应用图标 PNG
- `StartupWindowsVBS` — Windows 开机自启 VBS 脚本模板
- `StartupMacOSPlist` — macOS LaunchAgent plist 模板
- `StartupLinuxDesktop` — Linux .desktop 文件模板

### `info/version.go`

版本号定义，通过 `ldflags` 在构建时注入：

```bash
go build -ldflags="-X Doubao-Susurro/info.Version=v1.0.0" ./cmd/app
```

未注入时默认值为 `"dev"`。

## 依赖说明

| 依赖 | 用途 |
|------|------|
| `github.com/atotto/clipboard` | 系统剪贴板读写 |
| `github.com/energye/systray` | 跨平台系统托盘 |
| `github.com/go-vgo/robotgo` | 模拟键盘输入（自动输入功能） |
| `github.com/gofiber/fiber/v3` | Web 框架 |
| `github.com/mitchellh/mapstructure` | 结构体映射 |
| `github.com/spf13/viper` | 配置管理 |

### `internal/system/tray.go`

系统托盘管理，基于 `energye/systray`。

托盘菜单：
- **自动输入** — 复选框，切换 `auto_type` 配置并持久化
- **设置** — 打开浏览器配置页面
- **退出** — 退出程序

### `internal/system/lock/`

跨平台单实例进程锁。

- `lock.go` — `TryLock(name)` 在系统临时目录创建锁文件，返回 `unlock` 函数
- `lock_windows.go` — Windows 实现，使用 `kernel32.dll` 的 `LockFileEx` 排他锁
- `lock_unix.go` — Unix 占位实现（当前为空）

### `internal/tool/openbrowser.go`

- `OpenBrowser()` — 跨平台（Windows/macOS/Linux）打开默认浏览器访问 `http://localhost:{port}`

### `internal/tool/pngtoico.go`

- `PngToIco(pngData)` — 将 PNG 字节转换为 ICO 格式，用于系统托盘图标

### `internal/web/web.go`

基于 [Fiber v3](https://github.com/gofiber/fiber) 的 Web 服务，仅监听 `127.0.0.1`（禁止远程访问）。HTML 和图标通过 `//go:embed` 内嵌到二进制文件。

| 路由 | 方法 | 说明 |
|------|------|------|
| `/` | GET | Web 设置界面 |
| `/logo.png` | GET | 应用图标 |
| `/api/version` | GET | 获取版本号 |
| `/api/session` | GET | 获取当前 session |
| `/api/session/save` | POST | 保存 session（JSON body: `{"content": "..."}`) |
| `/api/poll` | GET | 手动触发一次消息获取（用于测试） |

### `info/version.go`

版本号定义，默认值为 `"dev"`，通过 `go build -ldflags` 在构建时注入。


## 已放弃的开发（如果有人想尝试也OK）

- 原始的 fyne 库实现：虽然现在 `github.com/fyne-io/systray` （或者其他从 `github.com/getlantern/systray` 衍生出的托盘界面）也能用，但是不支持高 DPI ，界面会很模糊，所以一直想用原始的 github.com/fyne-io/fyne 库实现，但是甚至没法运行他的 demo 。
- 注册为服务：尝试过 `github.com/kardianos/service` 库来注册为系统服务，结果可以 install / uninstall 但是无法启动。
