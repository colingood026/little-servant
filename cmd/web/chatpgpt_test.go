package main

import (
	"log"
	"testing"
)

func TestOpenAI_ChatWithChatGPT(t *testing.T) {

	ai := NewOpenAI("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")

	res, err := ai.ChatWithChatGPT("你好")

	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

}
