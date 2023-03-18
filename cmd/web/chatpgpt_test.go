package main

import (
	"log"
	"testing"
)

func TestOpenAI_ChatWithChatGPT(t *testing.T) {

	key := "TODO"
	ai := NewOpenAI(key)

	res, err := ai.ChatWithChatGPT("你好")

	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

}
