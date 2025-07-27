package models

type Question struct {
	ID            int    `json:"id"`
	QuestionText  string `json:"question_text"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
	OptionE       string `json:"option_e"`
	CorrectOption string `json:"correct_option"`
	ExamID        int    `json:"exam_id"`    // Untuk satu paket ujian tertentu
	SubjectID     int    `json:"subject_id"` // Tambahkan ini
	GradeID       int    `json:"grade_id"`   // Tambahkan ini
}
