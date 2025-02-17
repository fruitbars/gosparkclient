package main

import (
	"context"
	"fmt"
	"github.com/fruitbars/gosparkclient"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 从环境变量加载配置
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appID := os.Getenv("SPARKAI_APP_ID")
	apiKey := os.Getenv("SPARKAI_API_KEY")
	apiSecret := os.Getenv("SPARKAI_API_SECRET")
	hostURL := os.Getenv("SPARKAI_URL")
	domain := os.Getenv("SPARKAI_DOMAIN")

	// 创建客户端
	client, err := gosparkclient.NewSparkClient(
		gosparkclient.WithCredentials(appID, apiKey, apiSecret),
		gosparkclient.WithURLs(hostURL, ""), // embedding URL 可选
		gosparkclient.WithDomain(domain),
		gosparkclient.WithTimeout(time.Second*60),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// 示例4: 流式对话
	callbackChat(ctx, client)
}

func callbackChat(ctx context.Context, client *gosparkclient.SparkClient) {
	req := &gosparkclient.SparkChatRequest{
		Messages: []gosparkclient.SparkMessage{
			{
				Role:    "user",
				Content: "翻译为英文：你好",
			},
		},
		Temperature: 0.5,
		MaxTokens:   4096,
		TopK:        5,
	}

	// 定义回调函数来处理流式响应
	callback := func(resp *gosparkclient.SparkAPIResponse) {
		if len(resp.Payload.Choices.Text) > 0 {
			content := resp.Payload.Choices.Text[0].Content
			fmt.Print(content)
		}
	}

	// 使用回调方式调用
	if err := client.ChatWithCallback(ctx, req, callback); err != nil {
		log.Printf("Chat with callback failed: %v\n", err)
		return
	}

	fmt.Println()
}
