# Doubao Input

因为豆包的语音识别太准确有快速了，而且我想要手机上说话，电脑上就能直接粘贴，所以写了这个小工具。

他帮我在 deepseek 或者 github copilot 上输入我要跟他们描述的内容——而且比他们自带的语音输入准确很多，还支持中英文，还支持标点符号，支持很轻声音的说话，支持在嘈杂环境中说话，支持在有一两个别人在说话时也能准确的只把我的内容输出出来，他实在太好用了！

另外我还发现了一个妙用：我家里人要用电脑却不会电脑打字，用这个直接在豆包里说话再按下Ctrl-V就能输入了，太棒了！

所以我做了这个工具，希望也能帮助到有类似需求的朋友们。


## 功能特性

- 🎤 **语音转文字** — 利用豆包强大的语音识别，手机说话、电脑粘贴
- 🖥️ **系统托盘** — 最小化到托盘，不占用任务栏
- 🌐 **Web 设置界面** — 通过浏览器可视化配置 session
- 📋 **自动复制** — 新消息自动复制到剪贴板，直接 Ctrl+V 粘贴
- 🔄 **自动轮询** — 每秒检查新消息，实时同步
- ⚡ **开机自启** — 支持注册为系统服务，开机自动运行

## 下载

前往 [Releases](https://github.com/1299172402/Doubao-input/releases) 页面下载对应平台的可执行文件。

支持的平台：
- Windows (amd64)
- Linux (amd64)
- macOS (amd64 / arm64)

## Quick Start

### (首次使用) 1. 获取 session

1. 打开 [豆包网页版](https://www.doubao.com)，登录并进入一个对话
2. 按 `F12` 打开开发者工具 → **Network** / **网络** 标签
3. 在对话中发送一条消息，找到 `single` (`https://www.doubao.com/im/chain/single`) 请求
4. 右键该请求 → **Copy** / **复制** → **Copy as cURL (Bash)** / **复制为 cURL (Bash)**
5. 双击运行 `doubao-input.exe`
6. 在打开的浏览器页面中，按照提示粘贴从豆包复制的 cURL 内容
7. 点击「💾 保存配置」，然后点击「🚀 获取消息」测试是否正常

### 3. 在手机上对着豆包的同一个对话说话

### 4. Ctrl+V 粘贴到任何输入框中，享受语音输入的便利！

## Detail User Guide

### 系统托盘

系统托盘提供以下操作：

| 菜单项 | 说明 |
|--------|------|
| 打开设置页面 | 启动 Web 服务并自动打开浏览器 |
| 关闭设置页面 | 停止 Web 服务 |
| 开机自启 | 注册/注销 Windows 服务，实现开机自动运行(需要以管理员身份注册) |
| 退出 | 退出程序 |

### 命令行参数

| 参数 | 说明 |
|------|------|
| `-silent` | 静默模式，不打开浏览器和 Web 服务，仅运行后台轮询 |

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `DOUBAO_INPUT_PORT` | Web 服务端口 | `2828` |

## For Developers

### 构建并运行

```bash
go mod tidy
go build -o doubao-input.exe ./cmd
./doubao-input.exe
```

### 运行

```bash
go run ./cmd
```

### 项目结构

```
├── cmd/
│   ├── main.go           # 主程序入口，启动轮询和托盘
│   ├── tray.go           # 系统托盘菜单管理
│   ├── web.go            # Web 服务（Fiber 框架），提供配置界面和 API
│   ├── clipboard.go      # 消息轮询主循环，自动复制到剪贴板
│   ├── listener.go       # 调用豆包接口获取最新消息
│   ├── curl_parser.go    # cURL 命令解析与配置读写
│   ├── service.go        # Windows 服务管理（开机自启）
│   ├── tool.go           # 工具函数（PNG 转 ICO、打开浏览器）
│   └── static/
│       ├── index.html    # Web 设置界面
│       └── logo.png      # 应用图标
├── session.txt           # 存放从浏览器复制的 cURL 命令
├── go.mod                # Go 模块定义
└── README.md
```

### 各模块说明

#### `cmd/main.go`

主程序入口。支持 `-silent` 参数，启动消息轮询（后台）和系统托盘（阻塞主线程）。非静默模式下同时启动 Web 服务并自动打开浏览器。

#### `cmd/tray.go`

系统托盘管理，使用 `fyne.io/systray` 实现托盘图标和菜单。支持打开/关闭设置页面、开机自启开关、退出程序。

#### `cmd/web.go`

Web 服务，基于 `gofiber/fiber/v3`。HTML 和图标通过 `//go:embed` 内嵌到二进制文件中，无需额外资源文件。

| 路由 | 方法 | 说明 |
|------|------|------|
| `/` | GET | Web 设置界面 |
| `/logo.png` | GET | 应用图标 |
| `/api/version` | GET | 获取版本号（构建时注入） |
| `/api/session` | GET | 获取当前 session |
| `/api/session/save` | POST | 保存 session |
| `/api/poll` | GET | 手动获取最新消息（用于测试） |

#### `cmd/clipboard.go`

消息轮询主循环，每秒检查一次新消息。检测到新消息时自动写入系统剪贴板，并在控制台打印带时间戳的日志。

#### `cmd/listener.go`

- `GetLatestMessage(config)` — 根据 cURL 配置调用豆包接口，解析响应并返回最新一条用户消息（`user_type == 1`）

#### `cmd/curl_parser.go`

- `ReadCurlFile(path)` — 读取 cURL 文件内容，处理行尾反斜杠续行
- `ParseCurl(curlStr)` — 解析 cURL 命令，提取 URL、请求参数、请求头、Cookie、请求体
- `GetConfig(filePath)` — 读取文件并解析为 `CurlConfig` 结构
- `SaveCurlFile(path, content)` — 保存 cURL 内容到文件

#### `cmd/service.go`

Windows 服务管理，使用 `kardianos/service` 实现：
- `installService()` — 安装服务（开机自启）
- `uninstallService()` — 卸载服务
- `isServiceInstalled()` — 检查服务是否已安装

#### `cmd/tool.go`

- `pngToICO(pngData)` — 将 PNG 转换为 ICO 格式（用于系统托盘图标）
- `openBrowser(url)` — 跨平台打开默认浏览器

### 依赖

| 包名 | 用途 |
|------|------|
| `fyne.io/systray` | 系统托盘 |
| `github.com/atotto/clipboard` | 剪贴板操作 |
| `github.com/gofiber/fiber/v3` | Web 框架 |
| `github.com/kardianos/service` | Windows 服务管理 |
