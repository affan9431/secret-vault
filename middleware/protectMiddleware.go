package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware ensures only logged-in users can access the route
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fmt.Println(authHeader)
		if authHeader == "" {
			http.Error(w, "❌ Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Expecting: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // means no "Bearer " prefix
			http.Error(w, "❌ Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		secret := os.Getenv("JWT_SECRET_KEY")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "❌ Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Put token claims into request context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		http.Error(w, "❌ Could not read token claims", http.StatusUnauthorized)
	})
}
