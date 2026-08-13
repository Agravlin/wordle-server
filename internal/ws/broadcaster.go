package ws

import (
	"encoding/json"

	"github.com/agravlin/wordle-server/internal/game"
)

type WsMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func (r *Room) broadcastJSON(eventType string, payload any) {
	msg := WsMessage{
		Type:    eventType,
		Payload: payload,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		r.logger.Error("Failed to marshal broadcast message", "error", err)
		return
	}
	r.Broadcast <- bytes
}

func (r *Room) BroadcastFullState(gameState map[string]*game.Board) {
	r.broadcastJSON("FULL_STATE", gameState)
}

func (r *Room) BroadcastRowUpdate(nick string, board game.Board) {
	r.broadcastJSON("BOARD_UPDATE", map[string]any{
		"nick":  nick,
		"board": board,
	})
}
