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

func GetQuestions(w http.ResponseWriter, r *http.Request) {


	rows, err := config.DB.Query(
		"SELECT id, question_text, option_a, option_b, option_c, option_d, option_e, correct_option, exam_id FROM questions")
		if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
		}
		defer rows.Close()
		var question []models.Question

		for rows.Next() {
			var q models.Question
			err :=rows.Scan(&q.ID, &q.QuestionText, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD, &q.OptionE, &q.CorrectOption, &q.ExamID)
			if err != nil {
				http.Error(w, "Error scanning data", http.StatusInternalServerError)
				return
			}
			question = append(question, q)
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(question)
}

//tambah soal
func CreateQuestion(w http.ResponseWriter, r * http.Request) {
	// Cek apakah method yang digunakan POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Decode request body (JSON) ke dalam struct Question
	var q models.Question
	err := json.NewDecoder(r.Body).Decode(&q)

	fmt.Printf("DEBUG payload: %+v\n", q)

	
	if err != nil {
		http.Error(w, "Invalid JSON input"+err.Error(), http.StatusBadRequest)
		return
	}
	

	// Validasi CorrectOption
	validOptions := map[string]bool{"a": true, "b": true, "c": true, "d": true, "e": true}
	if _, ok := validOptions[q.CorrectOption]; !ok {
		http.Error(w, "Invalid correct_option, must be one of 'a', 'b', 'c', 'd', 'e'", http.StatusBadRequest)
		return
	}
	
	// Validasi question_text tidak kosong
if strings.TrimSpace(q.QuestionText) == "" {
    http.Error(w, "question_text cannot be empty", http.StatusBadRequest)
    return
}

// Validasi opsi jawaban tidak kosong
if strings.TrimSpace(q.OptionA) == "" ||
   strings.TrimSpace(q.OptionB) == "" ||
   strings.TrimSpace(q.OptionC) == "" ||
   strings.TrimSpace(q.OptionD) == "" ||
   strings.TrimSpace(q.OptionE) == "" {
    http.Error(w, "All options (a, b, c, d, e) must be provided and cannot be empty", http.StatusBadRequest)
    return
}
// Validasi exam_id (misal harus > 0)
	if q.ExamID <= 0 {
		http.Error(w, "exam_id must be provided and greater than 0", http.StatusBadRequest)
		return
	}

	// Siapkan statement untuk insert ke database
	stmt, err := config.DB.Prepare(
		`INSERT INTO questions
		(question_text, option_a, option_b, option_c, option_d, option_e, correct_option, exam_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
		if  err != nil {
			http.Error(w, "Database prepare failed: "+err.Error() ,http.StatusInternalServerError)
			return
		}

		defer stmt.Close()

		// Jalankan perintah insert
		fmt.Println("DEBUG: akan insert dengan exam_id =", q.ExamID)
		result, err :=stmt.Exec(q.QuestionText, q.OptionA, q.OptionB, q.OptionC, q.OptionD, q.OptionE, q.CorrectOption, q.ExamID)
		if err != nil {
			http.Error(w, "Insert failed"+err.Error(), http.StatusInternalServerError)
			return
		}
	
		lastInsertID, _ :=result.LastInsertId()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Question created",
			"id": lastInsertID,
		})
	}


//menampilkan per soal
func GetQuestionByID(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]

	var q models.Question
	err := config.DB.QueryRow(`SELECT id,	question_text, option_a, option_b, option_c, option_d, option_e, correct_option, exam_id FROM questions WHERE id = ?`, id).
	Scan(&q.ID, &q.QuestionText, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD, &q.OptionE, &q.CorrectOption, &q.ExamID)
	if err != nil {
		http.Error(w, "Question not found"+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(q)
}

//mengupdate per soal
func UpdateQuestionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var q models.Question
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		http.Error(w, "Invalid input"+err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := config.DB.Prepare(
		`UPDATE questions SET
		question_text = ?, option_a = ?, option_b = ?, option_c = ?, option_d = ?, option_e = ?, correct_option = ?, exam_id = ?
		WHERE id = ?`)
		if err != nil {
			http.Error(w, "Prepare failed"+err.Error(), http.StatusInternalServerError)
			return
		}

		defer stmt.Close()

		_, err = stmt.Exec(q.QuestionText, q.OptionA, q.OptionB, q.OptionC, q.OptionD, q.OptionE, q.CorrectOption, q.ExamID, id)
		if err != nil {
			http.Error(w, "Update failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Question updated"})
}

//delete soal
func DeleteQuestionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID"+err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := config.DB.Prepare("DELETE FROM questions WHERE id = ?")
	if err != nil {
		http.Error(w, "Prepare failed"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		http.Error(w, "Delete failed"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Question deleted"})
}