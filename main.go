package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"Chatapp/internal/chat"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("Username is required")
		return
	}
	// 1. EL SIKIŞMA (Handshake): TCP soketinin kontrolünü HTTP sunucusundan devral (Hijack).
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade hatası:", err)
		return
	}
	client := chat.NewClient(hub, conn, username)
	hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}	

func main() {
	hub := chat.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}