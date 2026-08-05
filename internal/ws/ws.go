package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all connections for now (CORS)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("New player connected!")

	// Simple echo loop for testing
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Connection closed or read error:", err)
			break
		}

		log.Printf("Received message: %s\n", p)

		// Echo the message back
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Message write error:", err)
			break
		}
	}
}
