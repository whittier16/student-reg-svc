package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/whittier16/student-reg-svc/internal/app/models"
	"github.com/whittier16/student-reg-svc/internal/pkg/database/db"
	"strings"
)

type RegisterRepository struct {
	DB *sql.DB
}

// NewRegisterRepository an instance of the RegisterRepository.
func NewRegisterRepository(db *db.MySQL) RegisterRepository {
	return RegisterRepository{DB: db.DBClient}
}

// Register sets the student and teacher emails in a new db record
func (rr *RegisterRepository) Register(ctx context.Context, teacherEmail string, studentEmails []string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegisterRepository.Register")
	defer span.Finish()

	var valueStrings []string
	var valueArgs []interface{}
	for _, email := range studentEmails {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, email)
		valueArgs = append(valueArgs, teacherEmail)
	}

	createRegisterStmt := "INSERT INTO register(student_id, teacher_id) VALUES %s"
	createRegisterStmt = fmt.Sprintf(createRegisterStmt, strings.Join(valueStrings, ", "))
	tx, err := rr.DB.Begin()
	if err != nil {
		log.Println("[Register][Register][Repository] Problem to querying to db, err: ", err.Error())
		return err
	}

	_, err = tx.Exec(createRegisterStmt, valueArgs...)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		log.Println("[Register][Register][Repository] Problem to querying to db, err: ", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return
}

// FindByEmailArr retrieves the emails with the given list of emails
func (rr *RegisterRepository) FindByEmailArr(ctx context.Context, emails []string) (resp []string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegisterRepository.FindByEmails")
	defer span.Finish()

	var valueStrings []string
	for _, email := range emails {
		valueStrings = append(valueStrings, fmt.Sprintf("'%s'", email))
	}

	var getRegQuery = "SELECT student_id FROM register WHERE suspended_on IS NULL AND teacher_id in (%s) GROUP BY student_id"
	getRegQuery = fmt.Sprintf(getRegQuery, strings.Join(valueStrings, ", "))

	rows, err := rr.DB.QueryContext(ctx, getRegQuery)
	if err != nil {
		log.Println("[Register][FindByEmailArr][Repository] Problem to querying to db, err: ", err.Error())
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	// An studentEmails slice to hold data from returned rows.
	var studentEmails []string

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var reg models.Register
		if err := rows.Scan(&reg.StudentID); err != nil {
			return studentEmails, err
		}
		studentEmails = append(studentEmails, reg.StudentID)
	}
	if err = rows.Err(); err != nil {
		log.Println("[Register][FindByEmailArr][Repository] Problem to querying to db, err: ", err.Error())
		return studentEmails, err
	}
	return studentEmails, nil
}

// Suspend sets the suspend date on an existing record on the requested email of the student
func (rr *RegisterRepository) Suspend(ctx context.Context, email string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegisterRepository.Update")
	defer span.Finish()

	_, err := rr.DB.ExecContext(ctx, updateRegisterQuery, email)
	if err != nil {
		log.Println("[Register][Suspend][Repository] Problem to querying to db, err: ", err.Error())
		return err
	}
	return err
}
