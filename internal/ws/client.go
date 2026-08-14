package ws

import (
	"log/slog"

	"github.com/agravlin/wordle-server/internal/game"
	"github.com/gorilla/websocket"
)

// Representing a single connected player
type Client struct {
	Nick string
	Room *Room
	Conn *websocket.Conn
	Send chan []byte
}

type ClientMessage struct {
	Action   Action `json:"action"`
	RowState []int  `json:"row_state,omitempty"`
	Guess    string `json:"guess,omitempty"`
}

type Action string

const (
	StateSync  Action = "SYNC_ROW"
	StateGuess Action = "GUESS"
	StateStart Action = "START_GAME"
)

func (c *Client) ReadPump() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg ClientMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			c.Room.logger.Error("Client disconnected or invalid JSON",
				slog.String("nick", c.Nick),
				slog.String("error", err.Error()),
			)
			break
		}

		c.handleMessage(msg)
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for {
		message, ok := <-c.Send
		if !ok {
			// Channel closed
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}

func (c *Client) handleMessage(msg ClientMessage) {
	switch msg.Action {
	case StateSync:
		c.handleSyncRow(msg.RowState)
	case StateGuess:
		c.handleGuess(msg.Guess)
	case StateStart:
		c.handleStartGame()
	default:
		c.Room.logger.Warn("Unknown action received",
			slog.String("action", string(msg.Action)),
		)
	}
}

func (c *Client) handleSyncRow(rowState []int) {
	if len(rowState) != 5 {
		c.Room.logger.Warn("Invalid row state length", slog.Int("length", len(rowState)))
		return
	}

	currentGame := c.Room.CurrentGame
	if currentGame == nil {
		return
	}

	board, exists := currentGame.GameState.Boards[c.Nick]
	if !exists || currentGame.State != game.StatePlaying {
		// Ignore if player is not in game or already finished
		return
	}

	c.Room.BroadcastRowUpdate(c.Nick, *board)
}

func (c *Client) handleGuess(guess string) {
	currentGame := c.Room.CurrentGame
	if currentGame == nil || currentGame.State != game.StatePlaying {
		return
	}

	_, err := currentGame.MakeGuess(c.Nick, guess)
	if err != nil {
		c.Room.logger.Warn("Guess failed",
			"nick", c.Nick,
			"error", err.Error(),
		)
		return
	}

	board, exists := currentGame.GameState.Boards[c.Nick]
	if !exists {
		return
	}

	c.Room.BroadcastRowUpdate(c.Nick, *board)

	if currentGame.State == game.StateFinished {
		c.Room.BroadcastFullState(currentGame.GameState.Boards)
	}
}

func (c *Client) handleStartGame() {
	var players []string
	for client := range c.Room.Clients {
		players = append(players, client.Nick)
	}

	c.Room.CurrentGame = game.NewGame(game.PickRandomWord(), players)

	c.Room.BroadcastFullState(c.Room.CurrentGame.GameState.Boards)
}
