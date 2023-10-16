package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/opentracing/opentracing-go"
	"github.com/whittier16/student-reg-svc/internal/app/models"
	"github.com/whittier16/student-reg-svc/internal/app/repository"
	"github.com/whittier16/student-reg-svc/internal/pkg/database/db"
	"strings"
)

// TODO implement interface
// Repo defines the DB level interaction of student registration
// type Repo interface {
//	Register(ctx context.Context, teacherEmail string, studentEmails []string) (err error)
// }

// Service uses repositories to provide an API for managing student registrations
type Service struct {
	sr repository.StudentRepository
	tr repository.TeacherRepository
	rr repository.RegisterRepository
}

// NewService returns a new instance of Service
func NewService(sr repository.StudentRepository, tr repository.TeacherRepository, rr repository.RegisterRepository) Service {
	return Service{
		sr: sr,
		tr: tr,
		rr: rr,
	}
}

// CreateStudent creates student record to the repo
func (s *Service) CreateStudent(ctx context.Context, params CreateStudentParams) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service.CreateStudent")
	defer span.Finish()

	if _, err := govalidator.ValidateStruct(params); err != nil {
		return err
	}

	tx, err := s.sr.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer a rollback in case anything fails.
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	entity := models.Student{
		Name:  params.Name,
		Email: params.Email,
	}
	err = s.sr.Create(ctx, &entity)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// GetStudent sends the request straight to the repo and retrieves student record
func (s *Service) GetStudent(ctx context.Context, email string) (models.Student, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service.GetStudent")
	defer span.Finish()

	student, err := s.sr.FindByEmail(ctx, email)
	switch {
	case err == nil:
	case errors.As(err, &db.ErrObjectNotFound{}):
		return models.Student{}, errors.New("student object not found")
	default:
		return models.Student{}, err
	}

	return student, nil
}

// GetTeacher sends the request straight to the repo and retrieves teacher record
func (s *Service) GetTeacher(ctx context.Context, email string) (models.Teacher, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service.GetTeacher")
	defer span.Finish()

	teacher, err := s.tr.FindByEmail(ctx, email)
	switch {
	case err == nil:
	case errors.As(err, &db.ErrObjectNotFound{}):
		return models.Teacher{}, errors.New("teacher object not found")
	default:
		return models.Teacher{}, err
	}

	return teacher, nil
}

// Register retrieves emails and passes of the created records to the repo
func (s *Service) Register(ctx context.Context, params RegisterStudentsParams) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service.Register")
	defer span.Finish()

	if _, err := govalidator.ValidateStruct(params); err != nil {
		return err
	}

	_, err := s.tr.FindByEmail(ctx, params.TeacherEmail)
	if err != nil {
		return err
	}

	res, err := s.sr.FindByEmailArr(ctx, params.StudentEmails, false)
	if len(res) < len(params.StudentEmails) {
		return errors.New("student not found in database")
	}

	tx, err := s.rr.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	err = s.rr.Register(ctx, params.TeacherEmail, params.StudentEmails)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// GetCommonStudents retrieves common students to the repo
func (s *Service) GetCommonStudents(ctx context.Context, params GetCommonStudentsParams) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service.GetCommonStudents")
	defer span.Finish()

	if _, err := govalidator.ValidateStruct(params); err != nil {
		return []string{}, err
	}

	cs, err := s.rr.FindByEmailArr(ctx, params.Teacher)
	switch {
	case err == nil:
	case errors.As(err, &db.ErrObjectNotFound{}):
		return []string{}, errors.New("teacher object not found")
	default:
		return []string{}, err
	}
	return cs, nil
}

// SendNotifications retrieves emails for sending notification
func (s *Service) SendNotifications(ctx context.Context, params SendNotificationsParams) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service.SendNotifications")
	defer span.Finish()

	if _, err := govalidator.ValidateStruct(params); err != nil {
		return []string{}, err
	}

	_, err := s.tr.FindByEmail(ctx, params.Teacher)
	if err != nil {
		return []string{}, err
	}

	teacherEmail := []string{}
	teacherEmail = append(teacherEmail, params.Teacher)
	resRegEmails, err := s.rr.FindByEmailArr(ctx, teacherEmail)
	if err != nil {
		return []string{}, err
	}

	txtNotif := params.Notifications
	txts := strings.Split(txtNotif, " @")
	notifEmails := []string{}
	for _, tx := range txts {
		i := strings.Index(tx, "@")
		if i > -1 {
			notifEmails = append(notifEmails, tx)
		}
	}

	resNotifEmails, err := s.sr.FindByEmailArr(ctx, notifEmails, false)
	if err != nil {
		return []string{}, err
	}

	resRegEmails = append(resRegEmails, resNotifEmails...)
	list := unique(resRegEmails)
	return list, err
}

// CreateTeacher creates teacher record to the repo
func (s *Service) CreateTeacher(ctx context.Context, params CreateTeacherParams) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service.CreateTeacher")
	defer span.Finish()

	if _, err := govalidator.ValidateStruct(params); err != nil {
		return err
	}

	tx, err := s.sr.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer a rollback in case anything fails.
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	entity := models.Teacher{
		Name:  params.Name,
		Email: params.Email,
	}
	err = s.tr.Create(ctx, &entity)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// Suspend retrieves student record and suspends a student
func (s *Service) Suspend(ctx context.Context, params SuspendStudentsParams) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service.Suspend")
	defer span.Finish()

	// find student object
	_, err := s.sr.FindByEmail(ctx, params.Student)
	if err != nil {
		return err
	}

	tx, err := s.sr.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	err = s.rr.Suspend(ctx, params.Student)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
