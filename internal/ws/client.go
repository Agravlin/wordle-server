package ws

import "github.com/gorilla/websocket"

// Representing a single connected player
type Client struct {
	room *Room
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for {
		message, ok := <-c.send
		if !ok {
			// Channel closed
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}
