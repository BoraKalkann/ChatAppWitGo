package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("EFN@#$1234567890_gizli_sirrimiz")
	}
	return []byte(secret)
}

func GenerateToken(username string) (string, error) {
claims := jwt.MapClaims{
	"username": username,
	"exp":    time.Now().Add(time.Hour * 24).Unix(),
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return token.SignedString(getJWTSecret())	
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJWTSecret(), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["username"].(string), nil
	}
	return "", fmt.Errorf("invalid token")
}
