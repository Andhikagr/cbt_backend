package handlers

import (
	"cbt-backend/config"
	"cbt-backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GET /grades
func GetGrades(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, name FROM grades")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var grades []models.Grade
	for rows.Next() {
		var g models.Grade
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		grades = append(grades, g)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

// POST /grades
func CreateGrade(w http.ResponseWriter, r *http.Request) {
	var g models.Grade
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec("INSERT INTO grades (name) VALUES (?)", g.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	g.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

// DELETE /grades/{id}
func DeleteGrade(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID"+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = config.DB.Exec("DELETE FROM grades WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//get gradesbyid
func GetGradeByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID"+err.Error(), http.StatusBadRequest)
        return
    }

    var g models.Grade
    err = config.DB.QueryRow("SELECT id, name FROM grades WHERE id = ?", id).Scan(&g.ID, &g.Name)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Grade not found"+err.Error(), http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(g)
}


