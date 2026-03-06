package chat

import "time"

type Message struct {
	Type      string    `bson:"type" json:"type"` 
	Sender    string    `bson:"sender" json:"sender"`
	Message   string    `bson:"message" json:"message"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

