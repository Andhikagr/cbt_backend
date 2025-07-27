package handlers

import (
	"cbt-backend/config"
	"cbt-backend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// --- Response Helpers ---
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// --- Validasi Soal ---
func validateQuestion(q models.Question) error {
	if strings.TrimSpace(q.QuestionText) == "" {
		return fmt.Errorf("question_text cannot be empty")
	}
	if strings.TrimSpace(q.OptionA) == "" ||
		strings.TrimSpace(q.OptionB) == "" ||
		strings.TrimSpace(q.OptionC) == "" ||
		strings.TrimSpace(q.OptionD) == "" ||
		strings.TrimSpace(q.OptionE) == "" {
		return fmt.Errorf("all options (a-e) must be provided and not empty")
	}
	validOptions := map[string]bool{"a": true, "b": true, "c": true, "d": true, "e": true}
	if !validOptions[q.CorrectOption] {
		return fmt.Errorf("correct_option must be one of a, b, c, d, or e")
	}
	if q.ExamID <= 0 || q.SubjectID <= 0 || q.GradeID <= 0 {
		return fmt.Errorf("exam_id, subject_id, and grade_id must be > 0")
	}
	return nil
}

// --- GET ALL Questions (with optional pagination) ---
func GetQuestions(w http.ResponseWriter, r *http.Request) {
	// Optional pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}
	offset := page * limit

	query := `SELECT id, question_text, option_a, option_b, option_c, option_d, option_e, correct_option, exam_id, subject_id, grade_id FROM questions LIMIT ? OFFSET ?`
	rows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error querying database")
		return
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		err := rows.Scan(
			&q.ID, &q.QuestionText, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD, &q.OptionE,
			&q.CorrectOption, &q.ExamID, &q.SubjectID, &q.GradeID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Error scanning data")
			return
		}
		questions = append(questions, q)
	}

	respondJSON(w, http.StatusOK, questions)
}

// --- GET Question by ID ---
func GetQuestionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var q models.Question
	err := config.DB.QueryRow(`SELECT id, question_text, option_a, option_b, option_c, option_d, option_e, correct_option, exam_id, subject_id, grade_id FROM questions WHERE id = ?`, id).
		Scan(&q.ID, &q.QuestionText, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD, &q.OptionE, &q.CorrectOption, &q.ExamID, &q.SubjectID, &q.GradeID)

	if err != nil {
		respondError(w, http.StatusNotFound, "Question not found")
		return
	}

	respondJSON(w, http.StatusOK, q)
}

// --- CREATE Question ---
func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var q models.Question
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON input: "+err.Error())
		return
	}

	if err := validateQuestion(q); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	stmt, err := config.DB.Prepare(`
		INSERT INTO questions
		(question_text, option_a, option_b, option_c, option_d, option_e, correct_option, exam_id, subject_id, grade_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database prepare failed: "+err.Error())
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(q.QuestionText, q.OptionA, q.OptionB, q.OptionC, q.OptionD, q.OptionE,
		q.CorrectOption, q.ExamID, q.SubjectID, q.GradeID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Insert failed: "+err.Error())
		return
	}

	lastInsertID, _ := result.LastInsertId()
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Question created",
		"id":      lastInsertID,
	})
}

// --- UPDATE Question by ID ---
func UpdateQuestionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var q models.Question
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	if err := validateQuestion(q); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	stmt, err := config.DB.Prepare(`
		UPDATE questions SET
		question_text = ?, option_a = ?, option_b = ?, option_c = ?, option_d = ?, option_e = ?, correct_option = ?, exam_id = ?, subject_id = ?, grade_id = ?
		WHERE id = ?`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Prepare failed: "+err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(q.QuestionText, q.OptionA, q.OptionB, q.OptionC, q.OptionD, q.OptionE,
		q.CorrectOption, q.ExamID, q.SubjectID, q.GradeID, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Update failed: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Question updated"})
}

// --- DELETE Question by ID ---
func DeleteQuestionByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID: "+err.Error())
		return
	}

	stmt, err := config.DB.Prepare("DELETE FROM questions WHERE id = ?")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Prepare failed: "+err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Delete failed: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Question deleted"})
}

// --- BULK INSERT Questions ---
func BulkInsertQuestions(w http.ResponseWriter, r *http.Request) {
	var questions []models.Question
	err := json.NewDecoder(r.Body).Decode(&questions)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	stmt, err := config.DB.Prepare(`
		INSERT INTO questions
		(question_text, option_a, option_b, option_c, option_d, option_e, correct_option, exam_id, subject_id, grade_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "DB prepare failed: "+err.Error())
		return
	}
	defer stmt.Close()

	var inserted int
	for _, q := range questions {
		if err := validateQuestion(q); err != nil {
			continue
		}
		if _, err := stmt.Exec(q.QuestionText, q.OptionA, q.OptionB, q.OptionC, q.OptionD, q.OptionE,
			q.CorrectOption, q.ExamID, q.SubjectID, q.GradeID); err == nil {
			inserted++
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":         "Bulk insert complete",
		"total_submitted": len(questions),
		"total_inserted":  inserted,
	})
}
