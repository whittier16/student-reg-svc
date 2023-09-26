package models

import "time"

type Student struct {
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	CreatedOn time.Time  `db:"created_on"`
	UpdatedOn *time.Time `db:"updated_on"`
	DeletedOn *time.Time `db:"deleted_on"`
}
