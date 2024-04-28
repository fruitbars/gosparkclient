# SparkClient Go Library
`SparkClient` 是一个Go语言库，用于与Spark AI的聊天API进行交互。它封装了创建请求、处理响应和WebSocket通信的逻辑，使得在Go应用程序中集成Spark AI服务变得简单。

## 安装
使用`go get`命令安装：
```bash
go get github.com/fruitbars/gosparkclient
```

## 快速开始
1. 确保你的项目中有一个`.env`文件，其中包含必要的环境变量，例如：
   ```
   SPARKAI_APP_ID=your_app_id
   SPARKAI_API_KEY=your_api_key
   SPARKAI_API_SECRET=your_api_secret
   SPARKAI_DOMAIN=your_domain
   SPARKAI_URL=your_base_url
   ```

2. 在你的Go代码中引入`SparkClient`库：
   ```go
   import "github.com/fruitbars/gosparkclient"
   ```

3. 使用`NewSparkClient`创建一个`SparkClient`实例，并使用其方法与API进行交互。

## 使用示例
以下是如何使用`SparkClient`发起一个简单的聊天请求的示例：

```go
package main

import (
    "log"

    "github.com/fruitbars/gosparkclient"
)

func main() {
    client := gosparkclient.NewSparkClient()

    // 发起一个简单的聊天请求
    r, sid, err := client.SparkChatSimple("你好")
    if err != nil {
        log.Fatalln(err)
    }
    log.Println("Session ID:", sid)
    log.Println("Response:", r)

    // 使用回调函数处理响应
    r, sid, err = client.SparkChatWithCallback(gosparkclient.SparkChatRequest{
        Prompt: "你好",
    }, func(response gosparkclient.SparkAPIResponse) {
        if len(response.Payload.Choices.Text) > 0 {
            log.Println("Callback Response:", response.Payload.Choices.Text[0].Content)
        }
    })
    if err != nil {
        log.Fatalln(err)
    }
    log.Println("Session ID:", sid)
    log.Println("Response:", r)
}
```


## 贡献
欢迎对`SparkClient`库做出贡献。如果你有任何问题或建议，请通过GitHub Issues提出。

## 许可证
`SparkClient`是在[MIT License](LICENSE)下发布的开源软件。