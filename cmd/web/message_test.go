package main

import (
	"fmt"
	"testing"
)

func TestInitMyMessage(t *testing.T) {
	//
	m := "單獨聊天我改成不用加[小僕人 ]開頭了"
	myMessage := InitMyMessage(m, false)
	fmt.Println(myMessage)
	myMessage = InitMyMessage(m, true)
	fmt.Println(myMessage)
}
