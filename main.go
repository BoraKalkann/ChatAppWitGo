package main

import (
	"log"
	"net/http"

	"Chatapp/config"
	"Chatapp/internal/chat"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(hub *chat.Hub, collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
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

	// Mevcut geçmiş mesajları, yeni bağlanan kullanıcıya gönder
	history := config.FetchHistory(collection)
	for _, msg := range history {
		client.Send(msg)
	}

	go client.WritePump()
	go client.ReadPump()
}	

func main() {
	collection := config.InitDB()
	hub := chat.NewHub(collection)

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, collection, w, r)
	})
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}