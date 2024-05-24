package main

import (
	"github.com/fruitbars/gosparkclient"
	"log"
)

func testDefeult() {
	client := gosparkclient.NewSparkClient()
	log.Println(client.AppID, client.HostURL)

	r, sid, err := client.SparkChatSimple("你好")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(sid, r)

	r, sid, err = client.SparkChatWithCallback(gosparkclient.SparkChatRequest{
		Message: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "user", Content: "Hello, how are you?"},
			{Role: "assistant", Content: "I'm fine, thank you!"},
		},
	}, func(response gosparkclient.SparkAPIResponse) {
		if len(response.Payload.Choices.Text) > 0 {
			log.Println(response.Header.Sid, response.Payload.Choices.Text[0].Content)
		}

	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(sid, r)
}

func testWithEnv(envName string, prompt string) {
	client := gosparkclient.NewSparkClientWithEnv(envName)
	log.Println(client.AppID, client.HostURL, client.Domain)

	r, sid, err := client.SparkChatWithCallback(gosparkclient.SparkChatRequest{
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

	log.Println(sid, r)
}

func main() {
	//testDefeult()
	testWithEnv("trans.env", "翻译为英文：你好")
}
