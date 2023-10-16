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

// Handler handles API requests
type Handler struct {
	logger *logrus.Logger
	router *mux.Router
	svc    services.Service
	cfg    *config.MainConfig
}

// New returns a new instance of the Handler
func New(log *logrus.Logger, db *db.MySQL, cfg *config.MainConfig) *Handler {
	return &Handler{
		logger: log,
		svc: services.NewService(
			repository.NewStudentRepository(db),
			repository.NewTeacherRepository(db),
			repository.NewRegisterRepository(db),
		),
		cfg: cfg,
	}
}

// Auth handles "POST /auth"
// Generates a new token.
// ---
// Responses:
//
//	204:
//	401:
func (h *Handler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := h.generateJWT()
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Token", token)
		w.Header().Set("Content-Type", "application/json")

		h.response(w, "", http.StatusNoContent)
	}
}

// Register handles "POST /api/register"
// Registers one or more students to a specified teacher.
// ---
// Responses:
//
//	204:
//	400:
//	401:
//	422:
func (h *Handler) Register() http.HandlerFunc {
	type request struct {
		TeacherEmail  string   `json:"teacher"`
		StudentEmails []string `json:"students"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := h.decode(r, &req)
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.svc.Register(r.Context(), services.RegisterStudentsParams{
			TeacherEmail:  req.TeacherEmail,
			StudentEmails: req.StudentEmails,
		})
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		h.response(w, "", http.StatusNoContent)
	}
}

// GetCommonStudents handles "GET /api/commonstudents"
// Gets list of common students to a given list of teachers.
// ---
// Responses:
//
//	200:
//	400:
//	401:
//	422:
func (h *Handler) GetCommonStudents() http.HandlerFunc {
	type response struct {
		Students []string `json:"students"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tq := r.URL.Query()["teacher"]
		if len(tq) == 0 {
			h.respondWithError(w, "missing required query params", http.StatusUnprocessableEntity)
			return
		}

		var emails []string
		for _, e := range tq {
			emails = append(emails, e)
		}
		res, err := h.svc.GetCommonStudents(r.Context(), services.GetCommonStudentsParams{
			Teacher: emails,
		})
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.response(w, response{
			Students: res,
		}, http.StatusOK)
	}
}

// Suspend handles "GET /api/suspend"
// Suspends a student.
// ---
// Responses:
//
//	204:
//	400:
//	401:
func (h *Handler) Suspend() http.HandlerFunc {
	type request struct {
		Student string `json:"student"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := h.decode(r, &req)
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.svc.Suspend(r.Context(), services.SuspendStudentsParams{
			Student: req.Student,
		})
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.response(w, "", http.StatusNoContent)
	}
}

// RetrieveNotifications handles "GET /api/retrievenotifications"
// Retrieves list of students who can receive a given notification.
// ---
// Responses:
//
//	200:
//	400:
//	401:
func (h *Handler) RetrieveNotifications() http.HandlerFunc {
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
		err := h.decode(r, &req)
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := h.svc.SendNotifications(r.Context(), services.SendNotificationsParams{
			Teacher:       req.Teacher,
			Notifications: req.Notification,
		})
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.response(w, response{
			Recipients: res,
		}, http.StatusOK)
	}
}

// CreateStudent handles "GET /api/student"
// Adds a student.
// ---
// Responses:
//
//	204:
//	400:
//	401:
func (h *Handler) CreateStudent() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := h.decode(r, &req)
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.svc.CreateStudent(r.Context(), services.CreateStudentParams{
			Email: req.Email,
			Name:  req.Name,
		})
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.response(w, "", http.StatusNoContent)
	}
}

// CreateTeacher handles "GET /api/teacher"
// Adds a teacher.
// ---
// Responses:
//
//	200:
//	400:
//	401:
func (h *Handler) CreateTeacher() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := h.decode(r, &req)
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.svc.CreateTeacher(r.Context(), services.CreateTeacherParams{
			Email: req.Email,
			Name:  req.Name,
		})
		if err != nil {
			h.respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.response(w, "", http.StatusNoContent)
	}
}
