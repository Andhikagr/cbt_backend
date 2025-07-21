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

// GET /answers
func GetAnswers(w http.ResponseWriter, r *http.Request) {
    rows, err := config.DB.Query("SELECT id, user_id, question_id, selected_option, is_correct, exam_id FROM answers")
    if err != nil {
        http.Error(w, "Failed to query answers", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var answers []models.Answer
    for rows.Next() {
        var a models.Answer
        if err := rows.Scan(&a.ID, &a.UserID, &a.QuestionID, &a.SelectedOption, &a.IsCorrect, &a.ExamID); err != nil {
            http.Error(w, "Failed to scan row", http.StatusInternalServerError)
            return
        }
        answers = append(answers, a)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(answers)
}
// GET /answers/{id}
func GetAnswerByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var a models.Answer
    err = config.DB.QueryRow(`
        SELECT id, user_id, question_id, selected_option, is_correct, exam_id 
        FROM answers WHERE id = ?`, id).
        Scan(&a.ID, &a.UserID, &a.QuestionID, &a.SelectedOption, &a.IsCorrect, &a.ExamID)
    
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Answer not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to query answer: "+err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(a)
}
