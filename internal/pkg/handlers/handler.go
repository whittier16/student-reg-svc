package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/whittier16/student-reg-svc/internal/app/config"
	"github.com/whittier16/student-reg-svc/internal/app/repository"
	"github.com/whittier16/student-reg-svc/internal/pkg/database/db"
	"github.com/whittier16/student-reg-svc/internal/pkg/services"
	"net/http"
)

type Service struct {
	logger *logrus.Logger
	router *mux.Router
	svc    services.Service
	cfg    *config.MainConfig
}

func New(lg *logrus.Logger, db *db.MySQL, cfg *config.MainConfig) *Service {
	return &Service{
		logger: lg,
		svc: services.NewService(
			repository.NewStudentRepository(db),
			repository.NewTeacherRepository(db),
			repository.NewRegisterRepository(db),
		),
		cfg: cfg,
	}
}

func (s *Service) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := s.generateJWT()
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Token", token)
		w.Header().Set("Content-Type", "application/json")

		s.response(w, "", http.StatusNoContent)
	}
}

func (s *Service) Register() http.HandlerFunc {
	type request struct {
		TeacherEmail  string   `json:"teacher"`
		StudentEmails []string `json:"students"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := s.decode(r, &req)
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.svc.Register(r.Context(), services.RegisterStudentsParams{
			TeacherEmail:  req.TeacherEmail,
			StudentEmails: req.StudentEmails,
		})
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.response(w, "", http.StatusNoContent)
	}
}

func (s *Service) GetCommonStudents() http.HandlerFunc {
	type response struct {
		Students []string `json:"students"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tq := r.URL.Query()["teacher"]
		if len(tq) == 0 {
			s.respondWithError(w, "missing required query params", http.StatusBadRequest)
			return
		}

		emails := []string{}
		for _, e := range tq {
			emails = append(emails, e)
		}
		res, err := s.svc.GetCommonStudents(r.Context(), services.GetCommonStudentsParams{
			Teacher: emails,
		})
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.response(w, response{
			Students: res,
		}, http.StatusOK)
	}
}

func (s *Service) Suspend() http.HandlerFunc {
	type request struct {
		Student string `json:"student"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := s.decode(r, &req)
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.svc.Suspend(r.Context(), services.SuspendStudentsParams{
			Student: req.Student,
		})
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.response(w, "", http.StatusNoContent)
	}
}

func (s *Service) RetrieveNotifications() http.HandlerFunc {
	type request struct {
		Teacher      string `json:"teacher"`
		Notification string `json:"notification"`
	}
	type response struct {
		Recipients []string `json:"recipients"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := s.decode(r, &req)
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := s.svc.SendNotifications(r.Context(), services.SendNotificationsParams{
			Teacher:       req.Teacher,
			Notifications: req.Notification,
		})
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.response(w, response{
			Recipients: res,
		}, http.StatusOK)
	}
}

func (s *Service) CreateStudent() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := s.decode(r, &req)
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.svc.CreateStudent(r.Context(), services.CreateStudentParams{
			Email: req.Email,
			Name:  req.Name,
		})
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.response(w, "", http.StatusNoContent)
	}
}

func (s *Service) CreateTeacher() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := s.decode(r, &req)
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.svc.CreateTeacher(r.Context(), services.CreateTeacherParams{
			Email: req.Email,
			Name:  req.Name,
		})
		if err != nil {
			s.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.response(w, "", http.StatusNoContent)
	}
}
