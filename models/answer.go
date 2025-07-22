package models

type Answer struct {
	ID             int     `json:"id"`
	UserID         *int    `json:"user_id,omitempty"`
	QuestionID     *int    `json:"question_id,omitempty"`
	SelectedOption *string `json:"selected_option,omitempty"`
	IsCorrect      *bool   `json:"is_correct,omitempty"`
	ExamID         int     `json:"exam_id"`
}
