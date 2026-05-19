# Doubao Input

因为豆包的语音识别太准确有快速了，而且我想要手机上说话，电脑上就能直接粘贴，所以写了这个小工具。

他帮我在 deepseek 或者 github copilot 上输入我要跟他们描述的内容——而且比他们自带的语音输入准确很多，还支持中英文，还支持标点符号，支持很轻声音的说话，支持在嘈杂环境中说话，支持在有一两个别人在说话时也能准确的只把我的内容输出出来，他实在太好用了！

另外我还发现了一个妙用：我家里人要用电脑却不会电脑打字，用这个直接在豆包里说话再按下Ctrl-V就能输入了，太棒了！

所以我做了这个工具，希望也能帮助到有类似需求的朋友们。


## 功能特性

- 🎤 **语音转文字** — 手机上对着豆包说话，电脑上自动获取识别结果
- 📋 **自动复制** — 新消息自动复制到剪贴板，直接 <kbd>Ctrl</kbd>+<kbd>V</kbd> 粘贴
- ⌨️ **自动输入** — 可选启用，自动将识别结果输入到当前焦点窗口
- 🖥️ **系统托盘** — 最小化到托盘，不占用任务栏，支持快速切换自动输入
- 🌐 **Web 设置界面** — 通过浏览器可视化配置 session
- 🔄 **实时轮询** — 每秒检查新消息，低延迟同步
- 🔒 **单实例运行** — 防止重复启动

## 下载

前往 [Releases](https://github.com/1299172402/Doubao-input/releases) 页面下载对应平台的可执行文件。

支持的平台：
- Windows (amd64)
- Linux (amd64)
- macOS (amd64 / arm64)

## Quick Start

### 1. 获取 Session（首次使用）

1. 打开 [豆包网页版](https://www.doubao.com)，登录并进入一个对话
2. 按 <kbd>F12</kbd> 打开开发者工具 → **Network**（网络）标签
3. 在对话中发送一条消息，找到 `single`（`https://www.doubao.com/im/chain/single`）请求
4. 右键该请求 → **Copy** → **Copy as cURL (Bash)**
5. 双击运行 `doubao-input.exe`，浏览器会自动打开设置页面
6. 将复制的 cURL 内容粘贴到文本框中
7. 点击「💾 保存配置」，然后点击「🚀 获取消息」测试是否正常

### 2. 在手机上对着豆包的同一个对话说话

### 3. Ctrl+V 粘贴到任何输入框中，享受语音输入的便利！

## Detail User Guide

### 系统托盘

系统托盘提供以下操作：

| 菜单项 | 说明 |
|--------|------|
| 自动输入 | 开启/关闭自动输入模式（自动将识别结果键入当前焦点窗口） |
| 设置 | 打开浏览器配置页面 |
| 退出 | 退出程序 |

### 命令行参数

| 参数 | 说明 |
|------|------|
| `-silent` | 静默模式，不自动打开浏览器（Web 服务和后台轮询仍正常运行） |

### 环境变量

环境变量前缀为 `DOUBAO_INPUT`，可覆盖配置文件中的对应字段：

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `DOUBAO_INPUT_PORT` | Web 服务端口 | `2828` |
| `DOUBAO_INPUT_AUTO_TYPE` | 是否启用自动输入 | `true` |


### 开机自启（Windows）

将以下内容保存为 `doubao-input-start.vbs`，放在与 `doubao-input.exe` 同目录下：

```vbs
Dim ws
Set ws = Wscript.CreateObject("Wscript.Shell")
ws.run "doubao-input-app.exe -silent",vbhide
Wscript.quit
```

然后为 `doubao-input-start.vbs` 创建快捷方式，将快捷方式放入 Windows 开始菜单启动文件夹：

```
%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup
```

之后每次开机都会自动静默启动程序。

