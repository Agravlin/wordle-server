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
	bytes, ok := r.jsonMessage(eventType, payload)
	if !ok {
		return
	}

	r.Broadcast <- bytes
}

func (r *Room) jsonMessage(eventType string, payload any) ([]byte, bool) {
	msg := WsMessage{
		Type:    eventType,
		Payload: payload,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		r.logger.Error("Failed to marshal broadcast message", "error", err)
		return nil, false
	}

	return bytes, true
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

func (r *Room) broadcastPlayerList() {
	var players []string
	for client := range r.Clients {
		players = append(players, client.Nick)
	}

	message, ok := r.jsonMessage("PLAYER_LIST", players)
	if ok {
		r.sendToClients(message)
	}
}
