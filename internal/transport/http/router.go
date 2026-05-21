package transporthttp

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/WebCraftersGH/Education-service/internal/middleware"
	swaggerdocs "github.com/WebCraftersGH/Education-service/internal/transport/http/docs"
	httphandlers "github.com/WebCraftersGH/Education-service/internal/transport/http/handlers"
	"github.com/WebCraftersGH/Education-service/pkg/logging"
)

func NewRouter(
	progressHandler *httphandlers.ProgressHandler,
	problemHandler *httphandlers.ProblemHandler,
	problemContentHandler *httphandlers.ProblemContentHandler,
	healthHandler *httphandlers.HealthHandler,
	docsHandler *swaggerdocs.DocsHandler,
	authChecker middleware.AuthChecker,
	logger logging.Logger,
	debugMode bool,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	//request-id-middleware
	router.Use(middleware.GenerateRequestID)
	router.Use(middleware.RequestLogger(logger))

	if debugMode {
		// DocsHandlers
		router.HandleFunc("/swagger/openapi.json", docsHandler.ServeSpec).Methods(http.MethodGet)
		router.HandleFunc("/swagger/", docsHandler.ServeUI).Methods(http.MethodGet)
		router.HandleFunc("/swagger", docsHandler.RedirectToUI).Methods(http.MethodGet)
	}

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
	api.HandleFunc("/problems/", problemHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/problems", problemHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/problems/{problemID}", problemHandler.ReadByID).Methods(http.MethodGet)
	api.HandleFunc("/problems/{slug}", problemHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/problems/{slug}", problemHandler.Delete).Methods(http.MethodDelete)

	// ProblemContentHandlers
	api.HandleFunc("/problems/{problemID}/content", problemContentHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/problems/{problemID}/content", problemContentHandler.ReadByProblemID).Methods(http.MethodGet)
	api.HandleFunc("/problems/{problemID}/content", problemContentHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/problems/{problemID}/content", problemContentHandler.DeleteByProblemID).Methods(http.MethodDelete)
	api.HandleFunc("/problems/{problemID}/content/solve", problemContentHandler.Solve).Methods(http.MethodPost)

	//Health
	api.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)

	return router
}
