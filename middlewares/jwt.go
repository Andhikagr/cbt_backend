package middlewares

import (
	"cbt-backend/config"
	"cbt-backend/models"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil token dari header Authorization: Bearer <token>
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &models.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			// ⬇️ Langsung ambil JWTSecret dari config, BUKAN var global
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Ambil role dari database berdasarkan user_id
		var role string
		err = config.DB.QueryRow("SELECT role FROM users WHERE id = ?", claims.ID).Scan(&role)
		if err != nil || role != "admin" {
			http.Error(w, "Only admin can access this resource", http.StatusForbidden)
			return
		}

		// Token valid dan user adalah admin
		next.ServeHTTP(w, r)
	})
}
