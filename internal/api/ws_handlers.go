package api

import (
	"log/slog"
	"net/http"

	"github.com/agravlin/wordle-server/internal/ws"
)

// GET /ws?room=abc&nick=alice
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	nick := r.URL.Query().Get("nick")

	room, err := h.svc.GetRoom(roomID)
	if err != nil {
		h.logger.Error("WS handshake failed",
			slog.String("error", err.Error()),
			slog.String("room_id", roomID),
			slog.String("nick", nick),
		)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if room.IsNickTaken(nick) {
		http.Error(w, "Nickname is already in use", http.StatusConflict)
		return
	}

	isHost, err := h.svc.IsHostJoining(roomID)
	if err != nil {
		http.Error(w, "Can not verify host status", http.StatusInternalServerError)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade failed",
			slog.String("error", err.Error()),
			slog.String("room_id", roomID),
			slog.String("nick", nick),
		)
		return
	}

	client := &ws.Client{
		Nick:   nick,
		Room:   room,
		IsHost: isHost,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	client.Room.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
