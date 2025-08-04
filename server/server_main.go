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
		OnConnect: func(clientID string) {
			fmt.Printf("Client connected: %s\n", clientID)
		},
		OnDisconnect: func(clientID string, err error) {
			fmt.Printf("Client disconnected: %s, error: %v\n", clientID, err)
		},
		OnMessage: func(clientID string, msg []byte) {
			fmt.Printf("Received message from client: %s, message: %s\n", clientID, string(msg))
		},
		OnError: func(err error) {
			fmt.Printf("Error: %v\n", err)
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
