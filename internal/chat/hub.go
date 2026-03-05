package chat

import "encoding/json"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			// Yeni kullanıcı bağlandığında herkese sistem mesajı gönder
			joinMsg := ClientMessage{
				Sender:  "SİSTEM",
				Message: "Sunucuya '" + client.username + "' olarak bağlanıldı.",
			}
			if data, err := json.Marshal(joinMsg); err == nil {
				for c := range h.clients {
					select {
					case c.send <- data:
					default:
						close(c.send)
						delete(h.clients, c)
					}
				}
			}

		case client := <-h.unregister:
			if h.clients[client] {
				delete(h.clients, client)
				close(client.send)

				// Kullanıcı ayrıldığında da sistem mesajı gönder (isteğe bağlı)
				leaveMsg := ClientMessage{
					Sender:  "SİSTEM",
					Message: " '" + client.username + "' sohbetten ayrıldı.",
				}
				if data, err := json.Marshal(leaveMsg); err == nil {
					for c := range h.clients {
						select {
						case c.send <- data:
						default:
							close(c.send)
							delete(h.clients, c)
						}
					}
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		}
	}
}