package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"Chatapp/internal/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitDB() *mongo.Collection {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB bağlantı hatası:", err)
	}

	Client = client

	
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB'ye erişilemiyor:", err)
	}


	return client.Database("chat_db").Collection("messages")
}



func FetchHistory(collection *mongo.Collection) [][]byte {
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	
	findOptions := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(50)
		
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Println("Geçmiş çekilirken hata:", err)
		return nil
	}
	defer cursor.Close(ctx)

	var messages []chat.Message
	if err = cursor.All(ctx, &messages); err != nil {
		log.Println("Cursor okuma hatası:", err)
		return nil
	}

	var history [][]byte
	for i := len(messages) - 1; i >= 0; i-- {
		msgBytes, _ := json.Marshal(messages[i])
		history = append(history, msgBytes)
	}

	return history
}