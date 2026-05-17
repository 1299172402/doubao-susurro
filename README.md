# Doubao Input

因为豆包的语音识别太准确有快速了，而且我想要手机上说话，电脑上就能直接粘贴，所以写了这个小工具。

他帮我在 deepseek 或者 github copilot 上输入我要跟他们描述的内容——而且比他们自带的语音输入准确很多，还支持中英文，还支持标点符号，支持很轻声音的说话，支持在嘈杂环境中说话，支持在有一两个别人在说话时也能准确的只把我的内容输出出来，他实在太好用了！

另外我还发现了一个妙用：我家里人要用电脑却不会电脑打字，用这个直接在豆包里说话再按下Ctrl-V就能输入了，太棒了！

所以我做了这个工具，希望也能帮助到有类似需求的朋友们。


## Quick Start

### 1. 安装依赖

```bash
pip install -r requirements.txt
```

### 2. 获取 session

1. 打开 [豆包网页版](https://www.doubao.com)，登录并进入一个对话
2. 按 `F12` 打开开发者工具 → **Network** / **网络** 标签
3. 在对话中发送一条消息，找到 `im/chain/single` 请求
4. 右键该请求 → **Copy** / **复制** → **Copy as cURL (Bash)** / **复制为 cURL (Bash)**
5. 将复制的内容粘贴到 `session.txt` 文件中

### 3. 验证配置

```bash
python curl_parser.py
```

正常输出 URL、PARAMS、HEADERS、COOKIES、PAYLOAD 即为成功。

### 4. 启动监听

```bash
python main.py
```

收到新消息时会自动：
- 在命令行打印消息内容和时间
- 将消息文本复制到剪贴板（`Ctrl+V` 粘贴）

按 `Ctrl+C` 停止。


## 项目结构

```
├── main.py          # 主程序入口，轮询消息并复制到剪贴板
├── listener.py      # 消息监听逻辑，调用接口获取最新用户消息
├── curl_parser.py   # curl 解析工具，从 session.txt 提取请求配置
├── session.txt      # 存放从浏览器复制的 curl 命令
└── requirements.txt # Python 依赖
```

## 各模块说明

### `session.txt`

存放从浏览器 DevTools 复制的原始 curl 命令。每次 Cookie 过期或切换会话时，重新复制粘贴即可。

### `curl_parser.py`

- `read_curl_file(path)` — 读取 curl 文件内容
- `parse_curl(curl_str)` — 解析 curl，提取 URL、请求参数、请求头、Cookie、请求体
- `get_config()` — 组装最终请求配置，修正 `direction=0`（向新消息方向拉取）

### `listener.py`

- `poll_new_message()` — 轮询一次接口，有新用户消息返回文本，无则返回 `None`

### `main.py`

主循环，每秒轮询一次，有新消息时打印并复制到剪贴板。
