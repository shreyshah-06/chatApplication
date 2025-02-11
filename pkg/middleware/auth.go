package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Get the secret key from environment variables
var secretKey = []byte(os.Getenv("SECRET_KEY"))

// JwtMiddleware validates the JWT token from the request
func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

		// Get the "Authorization" header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Missing Authorization Header")
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>" format
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			log.Println("Invalid Authorization Format")
			http.Error(w, "Invalid Authorization Format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		// Parse and validate the JWT token
		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Method)
				return nil, fmt.Errorf("unexpected signing method")
			}
			return secretKey, nil
		})

		if err != nil || !parsedToken.Valid {
			log.Printf("Invalid or expired token: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// If token is valid, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}