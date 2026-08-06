package ws

import "log/slog"

type Room struct {
	ID         string
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	TargetWord string
	logger     *slog.Logger
}

func NewRoom(id string, targetWord string, l *slog.Logger) *Room {
	r := &Room{
		ID:         id,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		TargetWord: targetWord,
		logger:     l.With(slog.String("room_id", id)), // Auto adds {"room_id": "AB12"} in each log
	}

	go r.Run()

	return r
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true
			r.logger.Info("New player joined the room", slog.String("nick", client.Nick))

		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
				r.logger.Info("Player left the room", slog.String("nick", client.Nick))
			}

		case message := <-r.Broadcast:
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(r.Clients, client)
				}
			}
		}
	}
}

func (r *Room) IsNickTaken(nick string) bool {
	for client := range r.Clients {
		if client.Nick == nick {
			return true
		}
	}
	return false
}
