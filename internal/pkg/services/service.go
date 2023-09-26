package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/whittier16/student-reg-svc/internal/app/models"
	"github.com/whittier16/student-reg-svc/internal/app/repository"
	"github.com/whittier16/student-reg-svc/internal/pkg/database/db"
	"strings"
)

type Service struct {
	sr repository.StudentRepository
	tr repository.TeacherRepository
	rr repository.RegisterRepository
}

func NewService(sr repository.StudentRepository, tr repository.TeacherRepository, rr repository.RegisterRepository) Service {
	return Service{
		sr: sr,
		tr: tr,
		rr: rr,
	}
}

func (s *Service) CreateStudent(ctx context.Context, params CreateStudentParams) error {
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

func (s *Service) GetStudent(ctx context.Context, email string) (models.Student, error) {
	todo, err := s.sr.FindByEmail(ctx, email)
	switch {
	case err == nil:
	case errors.As(err, &db.ErrObjectNotFound{}):
		return models.Student{}, errors.New("student object not found")
	default:
		return models.Student{}, err
	}
	return todo, nil
}

func (s *Service) GetTeacher(ctx context.Context, email string) (models.Teacher, error) {
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

func (s *Service) Register(ctx context.Context, params RegisterStudentsParams) error {
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

func (s *Service) GetCommonStudents(ctx context.Context, params GetCommonStudentsParams) ([]string, error) {
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

func (s *Service) SendNotifications(ctx context.Context, params SendNotificationsParams) ([]string, error) {
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

func (s *Service) CreateTeacher(ctx context.Context, params CreateTeacherParams) error {
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

func (s *Service) Suspend(ctx context.Context, params SuspendStudentsParams) error {
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
