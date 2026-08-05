package main

import (
	"log"
	"net/http"

	"github.com/agravlin/wordle-server/internal/ws"
)

func main() {
	http.HandleFunc("/ws", ws.HandleWebSocket)

	log.Println("WebSocket server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}
