package problemcatalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/WebCraftersGH/Education-service/internal/contracts"
	"github.com/WebCraftersGH/Education-service/internal/controller"
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProblemCatalogController struct {
	problemSVC        contracts.ProblemSVC
	problemContentSVC contracts.ProblemContentSVC
}

func NewProblemCatalogController(
	problemSVC contracts.ProblemSVC,
	problemContentSVC contracts.ProblemContentSVC,
) *ProblemCatalogController {
	return &ProblemCatalogController{
		problemSVC:        problemSVC,
		problemContentSVC: problemContentSVC,
	}
}

func (c *ProblemCatalogController) RegisterRoutes(r chi.Router) {
	r.Route("/problems", func(r chi.Router) {
		r.Post("/", c.CreateProblem)
		r.Get("/", c.ListProblems)
		r.Get("/{slug}", c.ReadProblemBySlug)
		r.Put("/{slug}", c.UpdateProblem)
		r.Delete("/{slug}", c.DeleteProblemBySlug)

		r.Post("/{problemID}/content", c.CreateProblemContent)
		r.Get("/{problemID}/content", c.ReadProblemContentByProblemID)
		r.Put("/{problemID}/content", c.UpdateProblemContent)
		r.Delete("/{problemID}/content", c.DeleteProblemContentByProblemID)
	})
}

func (c *ProblemCatalogController) CreateProblem(w http.ResponseWriter, r *http.Request) {
	var req createProblemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	problem := domain.Problem{
		Name:       strings.TrimSpace(req.Name),
		Slug:       strings.TrimSpace(req.Slug),
		Difficulty: strings.TrimSpace(req.Difficulty),
		Tag:        strings.TrimSpace(req.Tag),
		AuthorID:   req.AuthorID,
	}

	created, err := c.problemSVC.Create(r.Context(), problem)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toProblemResponse(created))
}

func (c *ProblemCatalogController) ReadProblemBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	problem, err := c.problemSVC.ReadBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "problem not found")
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toProblemResponse(problem))
}

func (c *ProblemCatalogController) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	current, err := c.problemSVC.ReadBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "problem not found")
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req updateProblemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedProblem := domain.Problem{
		ID:         current.ID,
		Name:       strings.TrimSpace(req.Name),
		Slug:       strings.TrimSpace(req.Slug),
		Difficulty: strings.TrimSpace(req.Difficulty),
		Tag:        strings.TrimSpace(req.Tag),
		Status:     current.Status,
		AuthorID:   req.AuthorID,
		VerifiedAt: current.VerifiedAt,
		CreatedAt:  current.CreatedAt,
		UpdatedAt:  current.UpdatedAt,
	}

	updated, err := c.problemSVC.Update(r.Context(), updatedProblem)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toProblemResponse(updated))
}

func (c *ProblemCatalogController) DeleteProblemBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if err := c.problemSVC.DeleteBySlug(r.Context(), slug); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "problem not found")
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *ProblemCatalogController) ListProblems(w http.ResponseWriter, r *http.Request) {
	filter := domain.ProblemFilter{
		Tag:        strings.TrimSpace(r.URL.Query().Get("tag")),
		Difficulty: strings.TrimSpace(r.URL.Query().Get("difficulty")),
		Limit:      parseIntQuery(r.URL.Query().Get("limit"), 20),
		Offset:     parseIntQuery(r.URL.Query().Get("offset"), 0),
	}

	problems, err := c.problemSVC.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toProblemsResponse(problems))
}

func (c *ProblemCatalogController) CreateProblemContent(w http.ResponseWriter, r *http.Request) {
	problemID, err := uuid.Parse(chi.URLParam(r, "problemID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	var req createProblemContentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	content := domain.ProblemContent{
		ProblemID:      problemID,
		DescriptionMD:  req.DescriptionMD,
		InputFormatMD:  req.InputFormatMD,
		OutputFormatMD: req.OutputFormatMD,
		ConstraintsMD:  req.ConstraintsMD,
		NotesMD:        req.NotesMD,
	}

	created, err := c.problemContentSVC.Create(r.Context(), content)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toProblemContentResponse(created))
}

func (c *ProblemCatalogController) ReadProblemContentByProblemID(w http.ResponseWriter, r *http.Request) {
	problemID, err := uuid.Parse(chi.URLParam(r, "problemID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	content, err := c.problemContentSVC.ReadByProblemID(r.Context(), problemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "problem content not found")
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toProblemContentResponse(content))
}

func (c *ProblemCatalogController) UpdateProblemContent(w http.ResponseWriter, r *http.Request) {
	problemID, err := uuid.Parse(chi.URLParam(r, "problemID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	current, err := c.problemContentSVC.ReadByProblemID(r.Context(), problemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "problem content not found")
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req updateProblemContentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedContent := domain.ProblemContent{
		ID:             current.ID,
		ProblemID:      current.ProblemID,
		DescriptionMD:  req.DescriptionMD,
		InputFormatMD:  req.InputFormatMD,
		OutputFormatMD: req.OutputFormatMD,
		ConstraintsMD:  req.ConstraintsMD,
		NotesMD:        req.NotesMD,
		CreatedAt:      current.CreatedAt,
		UpdatedAt:      current.UpdatedAt,
	}

	updated, err := c.problemContentSVC.Update(r.Context(), updatedContent)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toProblemContentResponse(updated))
}

func (c *ProblemCatalogController) DeleteProblemContentByProblemID(w http.ResponseWriter, r *http.Request) {
	problemID, err := uuid.Parse(chi.URLParam(r, "problemID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	if err := c.problemContentSVC.DeleteByProblemID(r.Context(), problemID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "problem content not found")
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, controller.ErrorResponse{Error: message})
}

func parseIntQuery(value string, def int) int {
	if strings.TrimSpace(value) == "" {
		return def
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return def
	}

	return n
}
