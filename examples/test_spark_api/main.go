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

	// 示例1: 简单对话
	fmt.Println("\n=== Simple Chat Example ===")
	simpleChat(ctx, client)

	/*
		// 示例2: 高级对话
		fmt.Println("\n=== Advanced Chat Example ===")
		advancedChat(ctx, client)

		// 示例3: 对话历史
		fmt.Println("\n=== Chat History Example ===")
		chatHistory(ctx, client)

	*/

	// 示例4: 流式对话
	callbackChat(ctx, client)
}

func simpleChat(ctx context.Context, client *gosparkclient.SparkClient) {
	resp, err := client.ChatSimple(ctx, "你好，请介绍一下你自己")
	if err != nil {
		log.Printf("Chat failed: %v\n", err)
		return
	}

	if len(resp.Payload.Choices.Text) > 0 {
		fmt.Printf("Assistant: %s\n", resp.Payload.Choices.Text[0].Content)
	}
}

func advancedChat(ctx context.Context, client *gosparkclient.SparkClient) {
	req := &gosparkclient.SparkChatRequest{
		Messages: []gosparkclient.SparkMessage{
			{
				Role:    "user",
				Content: "请帮我写一个Python函数，计算斐波那契数列的第n项",
			},
		},
		Temperature: 0.7,
		MaxTokens:   2000,
	}

	resp, err := client.Chat(ctx, req)
	if err != nil {
		log.Printf("Chat failed: %v\n", err)
		return
	}

	if len(resp.Payload.Choices.Text) > 0 {
		fmt.Printf("Assistant: %s\n", resp.Payload.Choices.Text[0].Content)
	}

	// 打印token使用情况
	fmt.Printf("\nToken Usage:\n")
	fmt.Printf("Question Tokens: %d\n", resp.Payload.Usage.Text.QuestionTokens)
	fmt.Printf("Completion Tokens: %d\n", resp.Payload.Usage.Text.CompletionTokens)
	fmt.Printf("Total Tokens: %d\n", resp.Payload.Usage.Text.TotalTokens)
}

func chatHistory(ctx context.Context, client *gosparkclient.SparkClient) {
	conversation := []gosparkclient.SparkMessage{
		{
			Role:    "user",
			Content: "你是谁？",
		},
		{
			Role:    "assistant",
			Content: "我是讯飞星火认知大模型，可以帮助你完成各种任务。",
		},
		{
			Role:    "user",
			Content: "你能做什么？",
		},
	}

	req := &gosparkclient.SparkChatRequest{
		Messages:    conversation,
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	resp, err := client.Chat(ctx, req)
	if err != nil {
		log.Printf("Chat failed: %v\n", err)
		return
	}

	if len(resp.Payload.Choices.Text) > 0 {
		fmt.Printf("Assistant: %s\n", resp.Payload.Choices.Text[0].Content)
	}
}

func callbackChat(ctx context.Context, client *gosparkclient.SparkClient) {
	req := &gosparkclient.SparkChatRequest{
		Messages: []gosparkclient.SparkMessage{
			{
				Role:    "user",
				Content: "请写一首诗歌，描写春天的美景",
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
			fmt.Println(content)
		}
	}

	// 使用回调方式调用
	if err := client.ChatWithCallback(ctx, req, callback); err != nil {
		log.Printf("Chat with callback failed: %v\n", err)
		return
	}
}
