package handlers

import (
	"cbt-backend/config"
	"cbt-backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// GET /users
func GetUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := config.DB.Query("SELECT id, name, email, role FROM users")
    if err != nil {
        http.Error(w, "Failed to query users", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role); err != nil {
            http.Error(w, "Failed to scan user", http.StatusInternalServerError)
            return
        }
        users = append(users, u)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// GET /users/{id}
func GetUserByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var u models.User
    err = config.DB.QueryRow("SELECT id, name, email, role FROM users WHERE id = ?", id).
        Scan(&u.ID, &u.Name, &u.Email, &u.Role)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "User not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to query user: "+err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(u)
}

// POST /users
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var u models.User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    if u.Name == "" || u.Email == "" || u.Password == "" {
        http.Error(w, "Name, email, and password are required", http.StatusBadRequest)
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }

    stmt, err := config.DB.Prepare("INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)")
    if err != nil {
        http.Error(w, "Prepare failed: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    result, err := stmt.Exec(u.Name, u.Email, string(hashedPassword), u.Role)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            http.Error(w, "Email already exists", http.StatusBadRequest)
            return
        }
        http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    lastID, _ := result.LastInsertId()
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "User created",
        "id":      lastID,
    })
}
// RegisterUser hanya boleh diakses oleh admin (role: "admin")
func RegisterUser(w http.ResponseWriter, r *http.Request) {
    // Middleware harus memastikan hanya admin yang bisa lewat
    var u models.User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    if u.Name == "" || u.Email == "" || u.Password == "" || u.Role == "" {
        http.Error(w, "Name, email, password, and role are required", http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }

    stmt, err := config.DB.Prepare("INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)")
    if err != nil {
        http.Error(w, "Prepare failed: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(u.Name, u.Email, string(hashedPassword), u.Role)
    if err != nil {
        http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "User registered successfully",
    })
}
// DELETE /users/{id}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    result, err := config.DB.Exec("DELETE FROM users WHERE id = ?", id)
    if err != nil {
        http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
        return
    }
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "User deleted successfully",
    })
}
