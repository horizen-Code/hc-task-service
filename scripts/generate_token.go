// scripts/generate_token.go
// Запуск: go run scripts/generate_token.go
package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Используем тот же секрет, что и в сервисе
	secret := "supersecretkey"
	userID := "22222222-2222-2222-2222-222222222222"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("=== JWT Token Generator ===")
	fmt.Printf("\nUser ID: %s\n", userID)
	fmt.Printf("Secret: %s\n", secret)
	fmt.Printf("Expires: %s\n\n", time.Now().Add(24*time.Hour).Format(time.RFC3339))
	fmt.Println("Token:")
	fmt.Println(tokenString)
	fmt.Println("\n=== Copy this token to Postman ===")
}
