package models

type Subject struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	GradeID *int   `json:"grade_id"` // gunakan *int agar bisa handle NULL juga
}
