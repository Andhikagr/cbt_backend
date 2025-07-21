package middlewares

import (
	"cbt-backend/models"
	"net/http"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userClaims, ok := r.Context().Value("user").(*models.JWTClaims)
		if !ok || userClaims.Role != "admin" {
			http.Error(w, "Admin only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
