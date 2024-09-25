package main

import (
	"encoding/json"
	"github.com/fruitbars/gosparkclient"
	"log"
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
		if response.Header.Code == 0 && len(response.Payload.Choices.Text) > 0 {
			log.Println(response.Payload.Choices.Text[0].Content)
		}

	})
	if err != nil {
		log.Fatalln(err, resp.Header.Sid)
	}

	log.Println(resp)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	testWithEnv(".env", "翻译为英文：你好")
}
