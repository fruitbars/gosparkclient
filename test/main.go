package main

import (
	"gosparkclient"
	"log"
)

func main() {
	client := gosparkclient.NewSparkClient()

	r, sid, err := client.SparkChatSimple("你好")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(sid, r)

	r, sid, err = client.SparkChatWithCallback(gosparkclient.SparkChatRequest{
		Prompt: "你好",
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
