package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/WebCraftersGH/Education-service/internal/apperrors"
	"github.com/WebCraftersGH/Education-service/internal/contracts"
	"github.com/WebCraftersGH/Education-service/internal/requestctx"
	"github.com/WebCraftersGH/Education-service/pkg/logging"
)

type ProgressHandler struct {
	usecase contracts.ProgressSVC
	logger  logging.Logger
}

func NewProgressHandler() *ProgressHandler {
	return &ProgressHandler{}
}

func (h *ProgressHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req SetProgressRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, ok := requestctx.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	created, err := h.usecase.SetProgress(r.Context(), userID, slug)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidArgument):
			writeError(w, http.StatusBadRequest, "invalid user id or slug")
			return
		case errors.Is(err, apperrors.ErrDuplicateRecord):
			writeError(w, http.StatusConflict, "progress already exists")
			return
		case errors.Is(err, apperrors.ErrInternal):
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		default:
			h.logger.WithError(err).Error("unexpected error in Progress Create handler")
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, toProgressResponse(created))
}

func (h *ProgressHandler) ListMyProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := requestctx.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	query := r.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	progressList, err := h.usecase.ProgressList(r.Context(), userID, limit, offset)

}
