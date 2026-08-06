package api

import (
	"log"
	"net/http"

	"github.com/agravlin/wordle-server/internal/ws"
)

// GET /ws?room=abc&nick=alice
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	nick := r.URL.Query().Get("nick")

	room, err := h.svc.GetRoom(roomID)
	if err != nil {
		log.Printf("WS Handshake Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v\n", err)
		return
	}

	client := &ws.Client{
		Nick: nick,
		Room: room,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	client.Room.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
