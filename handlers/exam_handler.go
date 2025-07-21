package handlers

import (
	"cbt-backend/config"
	"cbt-backend/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Create Exam
func CreateExam(w http.ResponseWriter, r *http.Request) {
	var exam models.Exam
	err := json.NewDecoder(r.Body).Decode(&exam)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := config.DB.Prepare(`
		INSERT INTO exams (subject_id, grade_id, start_time, end_time, duration_mins)
		VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		http.Error(w, "Prepare failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(exam.SubjectID, exam.GradeID, exam.StartTime, exam.EndTime, exam.DurationMins)
	if err != nil {
		http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Exam created successfully"})
}

// Get All Exams
func GetExams(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`SELECT id, subject_id, grade_id, start_time, end_time, duration_mins FROM exams`)
	if err != nil {
		http.Error(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var exams []models.Exam
	for rows.Next() {
		var e models.Exam
		err := rows.Scan(&e.ID, &e.SubjectID, &e.GradeID, &e.StartTime, &e.EndTime, &e.DurationMins)
		if err != nil {
			http.Error(w, "Scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		exams = append(exams, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exams)
}

// Get Exam by ID
func GetExamByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var e models.Exam

	err := config.DB.QueryRow(`
		SELECT id, subject_id, grade_id, start_time, end_time, duration_mins
		FROM exams WHERE id = ?`, id).
		Scan(&e.ID, &e.SubjectID, &e.GradeID, &e.StartTime, &e.EndTime, &e.DurationMins)

	if err != nil {
		http.Error(w, "Exam not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

// Update Exam
func UpdateExam(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var e models.Exam
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := config.DB.Prepare(`
		UPDATE exams SET subject_id = ?, grade_id = ?, start_time = ?, end_time = ?, duration_mins = ?
		WHERE id = ?`)
	if err != nil {
		http.Error(w, "Prepare failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.SubjectID, e.GradeID, e.StartTime, e.EndTime, e.DurationMins, id)
	if err != nil {
		http.Error(w, "Update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Exam updated successfully"})
}

// Delete Exam
func DeleteExam(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	stmt, err := config.DB.Prepare("DELETE FROM exams WHERE id = ?")
	if err != nil {
		http.Error(w, "Prepare failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		http.Error(w, "Delete failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Exam deleted successfully"})
}
