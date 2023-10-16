package repository

import (
	"context"
	"database/sql"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/whittier16/student-reg-svc/internal/app/models"
	"github.com/whittier16/student-reg-svc/internal/pkg/database/db"
)

type TeacherRepository struct {
	DB *sql.DB
}

// NewTeacherRepository an instance of the NewTeacherRepository.
func NewTeacherRepository(db *db.MySQL) TeacherRepository {
	return TeacherRepository{DB: db.DBClient}
}

// FindByEmail retrieves the teacher with the given email
func (tr *TeacherRepository) FindByEmail(ctx context.Context, email string) (resp models.Teacher, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TeacherRepository.FindByEmail")
	defer span.Finish()

	err = tr.DB.QueryRowContext(ctx, getTeacherQuery, email).Scan(
		&resp.Email,
		&resp.Name,
		&resp.CreatedOn,
		&resp.DeletedOn,
		&resp.UpdatedOn,
	)
	if err != nil {
		log.Println("[Teacher][FindByEmail][Repository] Problem to querying to db, err: ", err.Error())
		return resp, err
	}

	return resp, nil
}

// Create sets the teacher email and name in a new db record
func (tr *TeacherRepository) Create(ctx context.Context, input *models.Teacher) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TeacherRepository.Create")
	defer span.Finish()

	_, err = tr.DB.ExecContext(ctx, createTeacher,
		input.Email,
		input.Name,
	)
	if err != nil {
		log.Println("[Teacher][Create][Repository] Problem to querying to db, err: ", err.Error())
		return err
	}

	return nil
}
