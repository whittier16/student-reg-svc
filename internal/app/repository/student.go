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

type StudentRepository struct {
	DB *sql.DB
}

// NewStudentRepository an instance of the StudentRepository.
func NewStudentRepository(db *db.MySQL) StudentRepository {
	return StudentRepository{DB: db.DBClient}
}

// Create sets the email and name in a new db record
func (sr *StudentRepository) Create(ctx context.Context, input *models.Student) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "studentRepository.Create")
	defer span.Finish()

	_, err = sr.DB.ExecContext(ctx, createStudent,
		input.Email,
		input.Name,
	)
	if err != nil {
		log.Println("[Student][Create][Repository] Problem to querying to db, err: ", err.Error())
		return err
	}

	return nil
}

// FindByEmail retrieves the student with the given email
func (sr *StudentRepository) FindByEmail(ctx context.Context, email string) (resp models.Student, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StudentRepository.FindByEmail")
	defer span.Finish()

	err = sr.DB.QueryRowContext(ctx, getStudentQuery, email).Scan(
		&resp.Email,
		&resp.Name,
		&resp.CreatedOn,
		&resp.DeletedOn,
		&resp.UpdatedOn,
	)
	if err != nil {
		log.Println("[Student][FindByEmail][Repository] Problem to querying to db, err: ", err.Error())
		return resp, err
	}

	return resp, nil
}

// FindByEmailArr retrieves the emails with the given list of emails
func (sr *StudentRepository) FindByEmailArr(ctx context.Context, emails []string, isSuspended bool) (resp []string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StudentRepository.FindByEmailArr")
	defer span.Finish()

	var valueStrings []string
	for _, email := range emails {
		valueStrings = append(valueStrings, fmt.Sprintf("'%s'", email))
	}

	getQuery := "SELECT email FROM student JOIN register ON student.email = register.student_id WHERE student.email in (%s)%s"
	suspendedQry := ""
	if !isSuspended {
		suspendedQry = ` AND register.suspended_on IS NULL`
	}
	getQuery = fmt.Sprintf(getQuery, strings.Join(valueStrings, ", "), suspendedQry)
	fmt.Println(getQuery)
	rows, err := sr.DB.QueryContext(ctx, getQuery)
	if err != nil {
		log.Println("[Student][FindByEmailArr][Repository] Problem to querying to db, err: ", err.Error())
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("[Student][FindByEmailArr][Repository] Problem to querying to db, err: ", err.Error())
			return
		}
	}(rows)

	// studentEmails slice to hold data from returned rows.
	var studentEmails []string

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.Email); err != nil {
			return studentEmails, err
		}
		studentEmails = append(studentEmails, student.Email)
	}
	if err = rows.Err(); err != nil {
		log.Println("[Student][FindByEmailArr][Repository] Problem to querying to db, err: ", err.Error())
		return studentEmails, err
	}
	return studentEmails, nil
}
