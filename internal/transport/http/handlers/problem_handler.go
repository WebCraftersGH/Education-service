package handlers

import (
	"errors"
	"net/http"

	"github.com/WebCraftersGH/Education-service/internal/contracts"
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/WebCraftersGH/Education-service/internal/requestctx"
	"github.com/WebCraftersGH/Education-service/internal/slugify"
	"github.com/WebCraftersGH/Education-service/pkg/logging"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ProblemHandler struct {
	usecase contracts.ProblemSVC
	logger  logging.Logger
}

func NewProblemHandler(
	usecase contracts.ProblemSVC,
	logger logging.Logger,
) *ProblemHandler {
	return &ProblemHandler{usecase: usecase, logger: logger}
}

func (h *ProblemHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.WithFields(map[string]any{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Debug("create problem handler started")

	userID, ok := requestctx.UserID(r.Context())
	if !ok {
		h.logger.Debug("create problem rejected: missing user id in context")
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var createProblemRequest CreateProblemRequest
	if err := decodeJSON(r, &createProblemRequest); err != nil {
		h.logger.WithError(err).WithField("user_id", userID).Info("invalid create problem request body")
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	h.logger.WithFields(map[string]any{
		"user_id":    userID,
		"name":       createProblemRequest.Name,
		"difficulty": createProblemRequest.Difficulty,
		"tag":        createProblemRequest.Tag,
	}).Debug("create problem request decoded")

	domainProblem := domain.Problem{
		Name:       createProblemRequest.Name,
		Difficulty: createProblemRequest.Difficulty,
		Tag:        createProblemRequest.Tag,
		Slug:       slugify.Slugify(createProblemRequest.Name),
		Status:     domain.ProblemStatusDraft,
		AuthorID:   userID,
	}

	created, err := h.usecase.Create(r.Context(), domainProblem)
	if err != nil {
		h.logger.WithError(err).WithFields(map[string]any{
			"user_id":    userID,
			"name":       createProblemRequest.Name,
			"difficulty": createProblemRequest.Difficulty,
			"tag":        createProblemRequest.Tag,
		}).Debug("create problem usecase returned error")

		switch {
		case errors.Is(err, domain.ErrProblemNameRequired):
			writeError(w, http.StatusBadRequest, domain.ErrProblemNameRequired.Error())
			return
		case errors.Is(err, domain.ErrProblemSlugRequired):
			writeError(w, http.StatusBadRequest, domain.ErrProblemSlugRequired.Error())
			return
		case errors.Is(err, domain.ErrProblemDifficultyRequired):
			writeError(w, http.StatusBadRequest, domain.ErrProblemDifficultyRequired.Error())
			return
		case errors.Is(err, domain.ErrProblemAuthorIDRequired):
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		default:
			h.logger.WithError(err).Error("create problem error")
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	h.logger.WithFields(map[string]any{
		"id":        created.ID,
		"slug":      created.Slug,
		"author_id": created.AuthorID,
		"status":    created.Status,
	}).Debug("create problem handler completed")

	writeJSON(w, http.StatusCreated, toProblemResponse(created))
}

func (h *ProblemHandler) ReadByID(w http.ResponseWriter, r *http.Request) {
	rawProblemID := mux.Vars(r)["problemID"]
	h.logger.WithField("problem_id", rawProblemID).Debug("read problem by id handler started")

	problemID, err := uuid.Parse(rawProblemID)
	if err != nil {
		h.logger.WithError(err).WithField("problem_id", rawProblemID).Debug("invalid problem id")
		writeError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	problem, err := h.usecase.ReadByID(r.Context(), problemID)
	if err != nil {
		h.logger.WithError(err).WithField("problem_id", problemID).Debug("read problem by id usecase returned error")
		switch {
		case errors.Is(err, domain.ErrProblemIDRequired):
			writeError(w, http.StatusBadRequest, domain.ErrProblemIDRequired.Error())
		case errors.Is(err, gorm.ErrRecordNotFound):
			writeError(w, http.StatusNotFound, "problem not found")
		default:
			h.logger.WithError(err).Error("read problem by id error")
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	h.logger.WithFields(map[string]any{
		"id":   problem.ID,
		"slug": problem.Slug,
	}).Debug("read problem by id handler completed")

	writeJSON(w, http.StatusOK, toProblemResponse(problem))
}

func (h *ProblemHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := domain.ProblemFilter{
		Tag:        query.Get("tag"),
		Difficulty: query.Get("difficulty"),
		Limit:      parseIntQuery(query.Get("limit"), 20),
		Offset:     parseIntQuery(query.Get("offset"), 0),
	}

	h.logger.WithFields(map[string]any{
		"tag":        filter.Tag,
		"difficulty": filter.Difficulty,
		"limit":      filter.Limit,
		"offset":     filter.Offset,
	}).Debug("list problems handler started")

	problems, err := h.usecase.List(r.Context(), filter)
	if err != nil {
		h.logger.WithError(err).Error("list problems error")
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.logger.WithField("count", len(problems)).Debug("list problems handler completed")
	writeJSON(w, http.StatusOK, toProblemResponseList(problems))
}

func (h *ProblemHandler) Update(w http.ResponseWriter, r *http.Request) {}

func (h *ProblemHandler) Delete(w http.ResponseWriter, r *http.Request) {}
