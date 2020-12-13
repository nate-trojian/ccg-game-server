package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var _ = websocket.Upgrader{}

func main() {
	fmt.Println("just a test")
}