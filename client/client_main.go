package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"websocket/utils"
)

func main() {
	logger := log.New(os.Stdout, "[ws-client] ", log.LstdFlags|log.Llongfile)

	clientId := os.Args[1]

	config := NewClientConfig("ws", "localhost", "8080", "/ws", clientId, 5, 10)

	callbacks := &ClientCallbacks{
		Started: func() {
			fmt.Println("Client started")
		},
		Stopped: func() {
			fmt.Println("Client stopped")
		},
		OnConnect: func() {
			fmt.Println("Connected to server")
		},
		OnDisconnect: func(err error) {
			fmt.Println("Disconnected from server:", err)
		},
		OnMessage: func(msg []byte) {
			fmt.Println("Received message:", string(msg))
		},
		OnError: func(err error) {
			fmt.Println("Error:", err)
		},
	}

	client := NewClient(config, callbacks, logger)

	client.Start()

	go func() {
		time.Sleep(5 * time.Second)
		err := client.Send("Hello from Client " + clientId)
		if err != nil {
			logger.Println("Send failed:", err)
		}
	}()

	utils.CloseSignal()

	logger.Println("Shutting down client...")
	client.Stop()
}
