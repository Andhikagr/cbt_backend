package handlers

import (
	"cbt-backend/config"
	"cbt-backend/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds models.User

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Bersihkan email dari spasi
	creds.Email = strings.TrimSpace(creds.Email)

	var user models.User
	err := config.DB.QueryRow(`
		SELECT id, email, password AS password_hash, role 
		FROM users WHERE email = ?`, creds.Email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Role)

	if err != nil {
		log.Println("Login error:", err) // untuk debugging
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}


		/* 
		User mengisi email dan password dari aplikasi.
		Fungsi LoginUser:
		Mencocokkan email dengan database (users).
		Membandingkan password yang dikirim dengan hash di database.
		Kalau cocok → buat JWT token, isi:
		*/

	// Buat JWT token
	claims := models.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}
	/*
	Token ditandatangani dengan config.JWTSecret.
	Token dikirim ke client.
	*/


	// Kirim token ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}
