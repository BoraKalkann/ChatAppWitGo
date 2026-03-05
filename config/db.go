package config

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"Chatapp/internal/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitDB, MongoDB'ye bağlanır ve işlem yapılacak koleksiyonun referansını döner.
func InitDB() *mongo.Collection {
	// Bağlantı denemesi için 10 saniyelik bir zaman sınırı (timeout) belirliyoruz.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB bağlantı hatası:", err)
	}

	// Bağlantının fiziksel olarak ayakta olup olmadığını Ping ile doğruluyoruz.
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB'ye erişilemiyor:", err)
	}

	// chat_db veritabanındaki messages koleksiyonunun referansını döndürür.
	return client.Database("chat_db").Collection("messages")
}

// FetchHistory, veritabanından son mesajları okuyup,
// tarayıcının beklediği JSON []byte dilimine çevirir.
func FetchHistory(collection *mongo.Collection) [][]byte {
	// Veritabanı işlemleri için 5 saniyelik zaman aşımı (timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Sorgu ayarları: created_at alanına göre azalan (-1) sırala ve sadece 50 kayıt getir.
	findOptions := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(50)
	
	// Filtre boş (bson.M{}), yani tüm dökümanlar içinde bu ayarları uygula.
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

	// DİKKAT: Veriyi MongoDB'den en yeniden en eskiye doğru çektik (azalan sıralama yüzünden).
	// Tarayıcıda yukarıdan aşağıya (eskiden yeniye) doğru görünmesi için diziyi tersine çeviriyoruz.
	var history [][]byte
	for i := len(messages) - 1; i >= 0; i-- {
		// Her bir struct'ı, tarayıcının beklediği JSON formatına çevir
		msgBytes, _ := json.Marshal(messages[i])
		history = append(history, msgBytes)
	}

	return history
}