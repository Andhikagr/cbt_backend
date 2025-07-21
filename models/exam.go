package models

import "time"

type Exam struct {
	ID           int       `json:"id"`
	SubjectID    int       `json:"subject_id"`
	GradeID      int       `json:"grade_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	DurationMins int       `json:"duration_mins"`
}
