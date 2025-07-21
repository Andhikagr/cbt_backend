package models

type Answer struct {
	ID             int     `json:"id"`
	UserID         *int    `json:"user_id,omitempty"`         // nullable, pointer supaya bisa null
	QuestionID     *int    `json:"question_id,omitempty"`     // nullable
	SelectedOption *string `json:"selected_option,omitempty"` // nullable, nilai 'a', 'b', 'c', 'd', 'e'
	IsCorrect      *bool   `json:"is_correct,omitempty"`      // nullable, bisa true/false atau null
	ExamID         int     `json:"exam_id"`
}
