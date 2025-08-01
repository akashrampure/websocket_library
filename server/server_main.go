package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"websocket/utils"
)

func main() {
	logger := log.New(os.Stdout, "[ws-server] ", log.LstdFlags|log.Llongfile)

	config := NewWsConfig("localhost:8080", "/ws", []string{"*"})

	callbacks := &WsCallback{
		Started: func() {
			fmt.Println("Server started")
		},
		Stopped: func() {
			fmt.Println("Server stopped")
		},
		OnConnect: func() {
			fmt.Println("Client connected")
		},
		OnDisconnect: func(err error) {
			fmt.Println("Client disconnected", err)
		},
		OnMessage: func(msg []byte) {
			fmt.Println("Received message", string(msg))
		},
		OnError: func(err error) {
			fmt.Println("Error:", err)
		},
	}

	server := NewServer(config, callbacks, logger)

	go func() {
		if err := server.Start(); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	go func() {
		for {
			server.clients.Range(func(key, _ any) bool {
				clientID := key.(string)
				server.Send(clientID, "Hello Client "+clientID)
				return true
			})

			server.Broadcast("Hello All Clients!")
			time.Sleep(2 * time.Second)
		}
	}()

	utils.CloseSignal()

	logger.Println("Shutting down server...")
	server.Shutdown()
}
