package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/affan9431/secret-vault/models"
	"github.com/affan9431/secret-vault/storage"
	"github.com/affan9431/secret-vault/utils"
	"github.com/golang-jwt/jwt/v5"
)


func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var existingEmail string
	err = storage.DB.QueryRow("SELECT email FROM users WHERE email = ?", user.Email).Scan(&existingEmail)
	if err != sql.ErrNoRows {
		// If no error, that means email exists
		if err == nil {
			http.Error(w, "❌ Email already registered", http.StatusConflict)
			return
		}
		// Some other DB error
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		fmt.Println("Hash Password not created")
		return
	}

	_, err = storage.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "✅ User registered successfully!",
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	type LoginCredential struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var user LoginCredential
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var storedPassword string
	var userName string

	err = storage.DB.QueryRow("SELECT password, name FROM users WHERE email = ?", user.Email).Scan(&storedPassword, &userName)
	if err == sql.ErrNoRows {
		http.Error(w, "❌ Email not registered", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "❌ Server error", http.StatusInternalServerError)
		return
	}

	if !utils.CheckPasswordHash(user.Password, storedPassword) {
		http.Error(w, "❌ Invalid credentials!", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"email":    user.Email,
		"userName": userName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET_KEY")

	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		http.Error(w, "Could not sign token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+signedToken)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "✅ Login successfully!",
		"token":   signedToken,
	})

}
