package transporthttp

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/WebCraftersGH/Education-service/internal/middleware"
	swaggerdocs "github.com/WebCraftersGH/Education-service/internal/transport/http/docs"
	httphandlers "github.com/WebCraftersGH/Education-service/internal/transport/http/handlers"
)

func NewRouter(
	progressHandler *httphandlers.ProgressHandler,
	problemHandler *httphandlers.ProblemHandler,
	problemContentHandler *httphandlers.ProblemContentHandler,
	healthHandler *httphandlers.HealthHandler,
	docsHandler *swaggerdocs.DocsHandler,
	authChecker middleware.AuthChecker,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	//request-id-middleware
	router.Use(middleware.GenerateRequestID)

	// DocsHandlers
	router.HandleFunc("/swagger/openapi.json", docsHandler.ServeSpec).Methods(http.MethodGet)
	router.HandleFunc("/swagger/", docsHandler.ServeUI).Methods(http.MethodGet)
	router.HandleFunc("/swagger", docsHandler.RedirectToUI).Methods(http.MethodGet)

	//Health
	router.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)

	api := router.PathPrefix("/api/v1").Subrouter()

	//auth-middleware
	api.Use(middleware.AuthFromToken(authChecker))

	// ProgressHandlers
	api.HandleFunc("/me/progress", progressHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/me/progress", progressHandler.ListMyProgress).Methods(http.MethodGet)

	// ProblemHandlers
	api.HandleFunc("/problems", problemHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/problems", problemHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/problems/{slug}", problemHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/problems/{slug}", problemHandler.Delete).Methods(http.MethodDelete)

	//Health
	api.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)

	return router
}
