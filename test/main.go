package main

import (
	"encoding/json"
	"github.com/fruitbars/gosparkclient"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func testCallBack() {
	client := gosparkclient.NewSparkClient()
	log.Println(client.AppID, client.HostURL)
	resp, err := client.SparkChatWithCallback(gosparkclient.SparkChatRequest{
		Message: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "user", Content: "Hello, how are you?"},
		},
	}, func(response gosparkclient.SparkAPIResponse) {
		if len(response.Payload.Choices.Text) > 0 {
			//log.Println(response.Header.Sid, response.Payload.Choices.Text[0].Content)
			log.Println(response)
			// 将结构体转换为 JSON 字符串
			jsonData, err := json.Marshal(response)
			if err != nil {
				log.Println("Error marshalling to JSON:", err)
				return
			}

			// 输出 JSON 字符串
			log.Println(string(jsonData))
		}

	})
	if err != nil {
		log.Fatalln(err)
	}

	jsonData, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println("Error marshalling to JSON:", err)
		return
	}

	// 输出 JSON 字符串
	log.Println("all result:", string(jsonData))

}

func testDefeult() {
	client := gosparkclient.NewSparkClient()
	log.Println(client.AppID, client.HostURL)

	resp, err := client.SparkChatSimple("你好")
	if err != nil {
		log.Fatalln(err)
	}

	jsonData, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println("Error marshalling to JSON:", err)
		return
	}

	// 输出 JSON 字符串
	log.Println(string(jsonData))

}

func testWithEnv(envName string, prompt string) {
	client := gosparkclient.NewSparkClientWithEnv(envName)
	log.Println(client.AppID, client.HostURL, client.Domain)

	resp, err := client.SparkChatWithCallback(gosparkclient.SparkChatRequest{
		Message: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "user", Content: prompt},
		},
	}, func(response gosparkclient.SparkAPIResponse) {
		if len(response.Payload.Choices.Text) > 0 {
			log.Println(response.Header.Sid, response.Payload.Choices.Text[0].Content)
		}

	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(resp)
}

func testProxy() {
	proxyURL, err := url.Parse("http://127.0.0.1:8999")
	if err != nil {
		log.Fatalf("Failed to parse proxy URL: %v", err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(10) * time.Second,
			KeepAlive: time.Duration(10) * time.Second,
		}).DialContext,
	}

	client := gosparkclient.NewSparkClient()
	client.Transport = transport

	// 使用 SparkClient 发起请求
	response, err := client.SparkChatSimple("Hello!")
	if err != nil {
		log.Fatalf("SparkChatSimple failed: %v", err)
	}

	log.Printf("Response: %v", response)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//testDefeult()
	testWithEnv("spark-lite.env", "翻译为英文：你好")
	//testCallBack()
	//testProxy()
}
