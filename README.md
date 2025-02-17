# gosparkclient

[![Go Reference](https://pkg.go.dev/badge/github.com/fruitbars/gosparkclient.svg)](https://pkg.go.dev/github.com/yourusername/gosparkclient)
[![Go Report Card](https://goreportcard.com/badge/github.com/fruitbars/gosparkclient)](https://goreportcard.com/report/github.com/yourusername/gosparkclient)
[![License](https://img.shields.io/github/license/fruitbars/gosparkclient)](https://github.com/yourusername/gosparkclient/blob/main/LICENSE)

gosparkclient 是一个用 Go 语言编写的讯飞星火认知大模型 API 客户端库，提供了简洁、易用且功能完备的接口。

## 特性

- 支持全部星火认知大模型 API
- 支持流式输出（实时返回）
- 支持多种调用方式（标准、简单、回调）
- 优雅的错误处理
- 完整的类型定义
- 线程安全
- 详细的文档和示例

## 安装

```bash
go get github.com/yourusername/gosparkclient
```

## 快速开始

### 基础用法

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/yourusername/gosparkclient"
)

func main() {
    // 创建客户端
    client, err := gosparkclient.NewSparkClient(
        gosparkclient.WithCredentials("your-app-id", "your-api-key", "your-api-secret"),
        gosparkclient.WithURLs("wss://spark-api.xf-yun.com/v1.1/chat", ""),
        gosparkclient.WithDomain("generalv1.1"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 创建上下文
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    // 发送简单请求
    resp, err := client.ChatSimple(ctx, "你好，请介绍一下你自己")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp.Payload.Choices.Text[0].Content)
}
```

### 流式输出

```go
func streamingChat() {
    client, err := gosparkclient.NewSparkClient(...)
    if err != nil {
        log.Fatal(err)
    }

    req := &gosparkclient.SparkChatRequest{
        Messages: []gosparkclient.SparkMessage{
            {
                Role:    "user",
                Content: "请写一首诗歌",
            },
        },
    }

    callback := func(resp *gosparkclient.SparkAPIResponse) {
        if len(resp.Payload.Choices.Text) > 0 {
            fmt.Print(resp.Payload.Choices.Text[0].Content)
        }
    }

    err = client.ChatWithCallback(context.Background(), req, callback)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 更多参数设置

```go
func advancedChat() {
    client, err := gosparkclient.NewSparkClient(...)
    if err != nil {
        log.Fatal(err)
    }

    req := &gosparkclient.SparkChatRequest{
        Messages: []gosparkclient.SparkMessage{
            {
                Role:    "system",
                Content: "你是一个专业的程序员",
            },
            {
                Role:    "user",
                Content: "请写一个快速排序算法",
            },
        },
        Temperature: 0.5,
        MaxTokens:   4096,
    }

    resp, err := client.Chat(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp.Payload.Choices.Text[0].Content)
}
```

## 配置选项

支持以下配置选项：

```go
// 配置认证信息
WithCredentials(appID, apiKey, apiSecret string)

// 配置 API 地址
WithURLs(hostURL, embURL string)

// 配置模型域
WithDomain(domain string)

// 配置超时时间
WithTimeout(timeout time.Duration)

// 配置用户ID
WithUID(uid string)

// 配置审计选项
WithAuditing(auditing string)
```

## 错误处理

库提供了详细的错误类型：

- ConfigurationError: 配置错误
- ConnectionError: 连接错误
- AuthenticationError: 认证错误
- RequestError: 请求错误
- ResponseError: 响应错误
- WebSocketError: WebSocket 错误

每个错误都包含详细的错误信息和原始错误（如果有）。

## 示例

更多示例请查看 [examples](./examples) 目录。

## 贡献

欢迎提交 Issues 和 Pull Requests！

## 许可证

本项目采用 MIT 许可证，查看 [LICENSE](./LICENSE) 文件了解更多信息。

## 鸣谢

感谢所有贡献者以及以下开源项目：

- [gorilla/websocket](https://github.com/gorilla/websocket)
- [joho/godotenv](https://github.com/joho/godotenv)

## 相关文档

- [讯飞星火认知大模型 API 文档](https://www.xfyun.cn/doc/spark/overview.html)