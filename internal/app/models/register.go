package models

import (
	"time"
)

type Register struct {
	ID          int        `db:"id"`
	StudentID   string     `db:"student_id"`
	TeacherID   string     `db:"teacher_id"`
	CreatedOn   time.Time  `db:"created_on"`
	DeletedOn   *time.Time `db:"deleted_on"`
	SuspendedOn *time.Time `db:"suspended_on"`
}
