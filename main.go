package main

import (
	"log"
	"net/http"

	"Chatapp/config"
	"Chatapp/internal/chat"
	"Chatapp/internal/chat/auth"
	"Chatapp/internal/chat/upload"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
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
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		log.Println("Token is required")
		return
	}

	username, err := auth.ValidateToken(tokenString)
	if err != nil {
		log.Println("Geçersiz token:", err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade hatası:", err)
		return
	}
	client := chat.NewClient(hub, conn, username)
	hub.Register(client)

	history := config.FetchHistory(collection)
	for _, msg := range history {
		client.Send(msg)
	}

	go client.WritePump()
	go client.ReadPump()
}	

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Uyarı: .env dosyası bulunamadı, varsayılan değerler kullanılacak.")
	}

	collection := config.InitDB()
	hub := chat.NewHub(collection)

	go hub.Run()

	authCollection := config.Client.Database("test").Collection("users") 
	http.HandleFunc("/api/auth", auth.AuthHandler(authCollection))

	http.HandleFunc("/api/upload", upload.UploadHandler)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, collection, w, r)
	})

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}