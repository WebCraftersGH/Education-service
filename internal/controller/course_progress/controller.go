package courseprogress

import (
	"strings"
	"net/http"

	"github.com/go-chi/chi/v5"
	

	"github.com/WebCraftersGH/Education-service/internal/contracts"
	"github.com/WebCraftersGH/Education-service/internal/controller"
	"github.com/WebCraftersGH/Education-service/pkg/logging"
	"github.com/WebCraftersGH/Education-service/internal/requestctx"
)

type courseProgressController struct {
	svc contracts.ProgressSVC
	logger logging.Logger
}

func NewCourseProgressController(
	logger logging.Logger,
	svc contracts.ProgressSVC,
) *courseProgressController {
	return &courseProgressController{
		svc: svc,
		logger: logger,
	}
}

func (c *courseProgressController) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/progress", func(r chi.Router){
		r.Post("/", c.SetProgress)
		r.Get("/", c.ListProgressByUser)
	})
}

func (c *courseProgressController) SetProgress(w http.ResponseWriter, r *http.Request) {
	var req SetProgressRequest
	if err := controller.DecodeJSON(r, &req); err != nil {
		//TODO нуже лог?
		controller.WriteError(w, http.StatusBadRequest, "invalid request body")	
	}

	userID, ok := requestctx.UserID(r.Context())
	if !ok {
		controller.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		controller.WriteError(w, http.StatusBadRequest, "slug is nil")
	}

	created, err := c.svc.SetProgress(r.Context(), userID, slug)
	if err != nil {
		
	}

	controller.WriteJSON(w, http.StatusCreated, toProgressResponse(created))
}

func (c *courseProgressController) ListProgressByUser(w http.ResponseWriter, r *http.Request) {

}
