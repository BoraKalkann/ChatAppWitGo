package chat

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

type ClientMessage struct {
	Type      string    `json:"type" bson:"type"`
	Sender    string    `json:"sender" bson:"sender"`
	Message   string    `json:"message" bson:"message" `
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	username string
}

func (c *Client) Send(msg []byte) {
	c.send <- msg
}

func NewClient(hub *Hub, conn *websocket.Conn, username string) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		username: username,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)

	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		
		var msgObj ClientMessage
		err = json.Unmarshal(rawMessage, &msgObj)
		
		if err != nil || msgObj.Type == "" {
			msgObj = ClientMessage{
				Type:    "text",
				Sender:  c.username,
				Message: string(rawMessage),
			}
		} else {
			msgObj.Sender = c.username
		}

		jsonBytes, err := json.Marshal(msgObj)
		if err == nil {
			c.hub.broadcast <- jsonBytes
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
