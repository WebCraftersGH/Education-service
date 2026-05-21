package handlers

import (
	"errors"
	"net/http"

	"github.com/WebCraftersGH/Education-service/internal/contracts"
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/WebCraftersGH/Education-service/internal/requestctx"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ProblemContentHandler struct {
	usecase contracts.ProblemContentSVC
}

func NewProblemContentHandler(usecase contracts.ProblemContentSVC) *ProblemContentHandler {
	return &ProblemContentHandler{usecase: usecase}
}

func (h *ProblemContentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := requestctx.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	problemID, ok := parseProblemID(w, r)
	if !ok {
		return
	}

	var req CreateProblemContentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.usecase.Create(r.Context(), domain.ProblemContent{
		ProblemID:     problemID,
		AuthorID:      userID,
		ActualGraph:   req.ActualGraph,
		ExpectedGraph: req.ExpectedGraph,
		FullText:      req.FullText,
	})
	if err != nil {
		writeProblemContentError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toProblemContentResponse(created))
}

func (h *ProblemContentHandler) ReadByProblemID(w http.ResponseWriter, r *http.Request) {
	problemID, ok := parseProblemID(w, r)
	if !ok {
		return
	}

	content, err := h.usecase.ReadByProblemID(r.Context(), problemID)
	if err != nil {
		writeProblemContentError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toProblemContentResponse(content))
}

func (h *ProblemContentHandler) Update(w http.ResponseWriter, r *http.Request) {
	problemID, ok := parseProblemID(w, r)
	if !ok {
		return
	}

	current, err := h.usecase.ReadByProblemID(r.Context(), problemID)
	if err != nil {
		writeProblemContentError(w, err)
		return
	}

	var req UpdateProblemContentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.usecase.Update(r.Context(), domain.ProblemContent{
		ID:            current.ID,
		ProblemID:     problemID,
		AuthorID:      current.AuthorID,
		ActualGraph:   req.ActualGraph,
		ExpectedGraph: req.ExpectedGraph,
		FullText:      req.FullText,
	})
	if err != nil {
		writeProblemContentError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toProblemContentResponse(updated))
}

func (h *ProblemContentHandler) DeleteByProblemID(w http.ResponseWriter, r *http.Request) {
	problemID, ok := parseProblemID(w, r)
	if !ok {
		return
	}

	if err := h.usecase.DeleteByProblemID(r.Context(), problemID); err != nil {
		writeProblemContentError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProblemContentHandler) Solve(w http.ResponseWriter, r *http.Request) {
	problemID, ok := parseProblemID(w, r)
	if !ok {
		return
	}

	var req SolveProblemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	solved, err := h.usecase.Solve(r.Context(), problemID, req.ActualGraph)
	if err != nil {
		writeProblemContentError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, SolveProblemResponse{Solved: solved})
}

func parseProblemID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	rawProblemID := mux.Vars(r)["problemID"]
	problemID, err := uuid.Parse(rawProblemID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid problem id")
		return uuid.Nil, false
	}

	return problemID, true
}

func writeProblemContentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProblemContentProblemIDRequired):
		writeError(w, http.StatusBadRequest, domain.ErrProblemContentProblemIDRequired.Error())
	case errors.Is(err, domain.ErrProblemAuthorIDRequired):
		writeError(w, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, domain.ErrProblemContentIDRequired):
		writeError(w, http.StatusBadRequest, domain.ErrProblemContentIDRequired.Error())
	case errors.Is(err, domain.ErrProblemContentFullTextRequired):
		writeError(w, http.StatusBadRequest, domain.ErrProblemContentFullTextRequired.Error())
	case errors.Is(err, domain.ErrProblemContentActualGraphInvalid):
		writeError(w, http.StatusBadRequest, domain.ErrProblemContentActualGraphInvalid.Error())
	case errors.Is(err, domain.ErrProblemContentExpectedGraphInvalid):
		writeError(w, http.StatusBadRequest, domain.ErrProblemContentExpectedGraphInvalid.Error())
	case errors.Is(err, gorm.ErrRecordNotFound):
		writeError(w, http.StatusNotFound, "problem content not found")
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}
