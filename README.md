
# QQBot-Go 集成 Pixiv 图片搜索与 DeepSeek 聊天功能

这个项目是一个用 Go 实现的 QQ bot，集成了 Pixiv 图片搜索和 DeepSeek 聊天功能。该 bot 使用 webhook 接收消息，并通过 RPC 转发到本地进行处理。服务器端运行在远程服务器上，其余组件则在本机上运行。

## 🌟 核心功能
- 📡 Webhook实时消息接收
- 🔄 反向RPC消息转发
- 💓 TCP心跳检测（30秒间隔）
- 🚀 并发消息队列处理（Goroutine + Channel）
- 🧠 DeepSeek智能对话集成
- 🖼️ 图片搜索/生成功能
- ⚡ 客户端连接动态管理
- 🔄 自动断线重连机制
- 更多功能待开发...

## 项目配置

### 1. **服务器端配置（远程服务器）**

服务器端部分（如 webhook 处理、RPC 服务器、消息处理）需要运行在远程服务器上。

#### 环境要求
- 服务器上需要安装 Go（Golang）。
- 需要有一个运行中的 QQ bot 实例。
- 服务器需要能够访问互联网，以便与 QQ 和 DeepSeek 进行通信。

#### 运行服务器端步骤：
1. 克隆仓库到服务器：
   ```bash
   git clone https://github.com/your-repository/QQbot-go.git
   cd QQbot-go/QQBotWebHook
2. 配置服务器设置（AppID:"",AppSecret:"" ）
3. 运行代码：
   ```bash
   go mod tidy
   go run main.go
#### 客户端配置（本机）
客户端部分（如 Pixiv 搜索、DeepSeek 聊天和发送消息请求）需要在本机上运行
1. 克隆仓库到本地主机：
   ```bash
   git clone https://github.com/your-repository/QQbot-go.git
   cd QQbot-go
   rm -rf QQBotWebHook
2. 配置本机设置（Pixiv PHPSESSID、DeepSeek API_KYES）
3. ```bash
   go mod tidy
   go run main.go

