package handlers

import (
	"errors"
	"github.com/WebCraftersGH/Education-service/internal/contracts"
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/WebCraftersGH/Education-service/internal/slugify"
	"github.com/WebCraftersGH/Education-service/pkg/logging"
	"net/http"
)

type ProblemHandler struct {
	usecase contracts.ProblemSVC
	logger  logging.Logger
}

func NewProblemHandler(
	usecase contracts.ProblemSVC,
	logger logging.Logger,
) *ProblemHandler {
	return &ProblemHandler{usecase: usecase}
}

func (h *ProblemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createProblemRequest CreateProblemRequest
	if err := decodeJSON(r, &createProblemRequest); err != nil {
		h.logger.WithError(err).Info("invalid request body")
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	domainProblem := domain.Problem{
		Name:       createProblemRequest.Name,
		Difficulty: createProblemRequest.Difficulty,
		Tag:        createProblemRequest.Tag,
		Slug:       slugify.Slugify(createProblemRequest.Name),
		Status:     domain.ProblemStatusDraft,
	}

	created, err := h.usecase.Create(r.Context(), domainProblem)
	if err != nil {
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
		default:
			h.logger.WithError(err).Error("create problem error")
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *ProblemHandler) List(w http.ResponseWriter, r *http.Request) {}

func (h *ProblemHandler) Update(w http.ResponseWriter, r *http.Request) {}

func (h *ProblemHandler) Delete(w http.ResponseWriter, r *http.Request) {}
