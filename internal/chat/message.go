package chat

import "time"

// Message, sohbet mesajlarının hem MongoDB'de hem de
// WebSocket üzerinden taşınan ortak veri modelidir.
type Message struct {
	Sender    string    `bson:"sender" json:"sender"`
	Message   string    `bson:"message" json:"message"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

