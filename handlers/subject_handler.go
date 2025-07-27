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
)

// GET /subjects
func GetAllSubjects(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, name, grade_id FROM subjects")
	if err != nil {
		http.Error(w, "Failed to query subjects", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.ID, &s.Name, &s.GradeID); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		subjects = append(subjects, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}

// POST /subjects
func CreateSubject(w http.ResponseWriter, r *http.Request) {
	var s models.Subject
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(s.Name) == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	stmt, err := config.DB.Prepare("INSERT INTO subjects (name, grade_id) VALUES (?, ?)")
	if err != nil {
		http.Error(w, "Prepare failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(s.Name, s.GradeID)
	if err != nil {
		
		if strings.Contains(err.Error(), "Duplicate entry") {
			http.Error(w, "Subject with that name already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, _ := result.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Subject created",
		"id":      lastID,
	})
}

// DELETE /subjects/{id}
func DeleteSubject(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    result, err := config.DB.Exec("DELETE FROM subjects WHERE id = ?", id)
    if err != nil {
        http.Error(w, "Failed to delete subject: "+err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, "Failed to get affected rows: "+err.Error(), http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Subject not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent) 
}

// GET /subjects/{id}
func GetSubjectByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var s models.Subject
    err = config.DB.QueryRow("SELECT id, name, grade_id FROM subjects WHERE id = ?", id).
    Scan(&s.ID, &s.Name, &s.GradeID)

    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Subject not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to query subject: "+err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(s)
}

// GET /grades/{id}/subjects
func GetSubjectsByGradeID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gradeID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid grade ID", http.StatusBadRequest)
		return
	}

	rows, err := config.DB.Query("SELECT id, name, grade_id FROM subjects WHERE grade_id = ?", gradeID)
	if err != nil {
		http.Error(w, "Failed to query subjects by grade", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.ID, &s.Name, &s.GradeID); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		subjects = append(subjects, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}
