package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"Chatapp/internal/chat/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func AuthHandler(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Sadece POST metodu desteklenir", http.StatusMethodNotAllowed)
			return
		}

		var reqUser models.User
		if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
			http.Error(w, "Geçersiz veri", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var dbUser models.User
		err := collection.FindOne(ctx, bson.M{"username": reqUser.Username}).Decode(&dbUser)

		if err == mongo.ErrNoDocuments {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqUser.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Şifreleme hatası", http.StatusInternalServerError)
				return
			}
			reqUser.Password = string(hashedPassword)

			_, err = collection.InsertOne(ctx, reqUser)
			if err != nil {
				http.Error(w, "Kayıt hatası", http.StatusInternalServerError)
				return
			}
					
			token, err := GenerateToken(reqUser.Username)
			if err != nil {
				http.Error(w, "Token oluşturma hatası", http.StatusInternalServerError)
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"token": token, "message": "Kayıt başarılı"})
			return

		} else if err != nil {
			http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
			return
		}

		
		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(reqUser.Password))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Hatalı şifre"})
			return
		}
		
		token, err := GenerateToken(reqUser.Username)
		if err != nil {
			http.Error(w, "Token oluşturma hatası", http.StatusInternalServerError)
			return
		}

			
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token, "message": "Giriş başarılı"})
	}
}
