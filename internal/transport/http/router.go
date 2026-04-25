package transporthttp

import (
	"net/http"

	"github.com/gorilla/mux"

	swaggerdocs "github.com/WebCraftersGH/Education-service/internal/transport/http/docs"
	httphandlers "github.com/WebCraftersGH/Education-service/internal/transport/http/handlers"
)

const servicePrefix = "edu"

func NewRouter(
	progressHandler *httphandlers.ProgressHandler,
	problemHandler *httphandlers.ProblemHandler,
	docsHandler *swaggerdocs.DocsHandler,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	baseRouter := router.PathPrefix("/" + servicePrefix).Subrouter()

	// DocsHandlers
	baseRouter.HandleFunc("/swagger/openapi.json", docsHandler.ServeSpec).Methods(http.MethodGet)
	baseRouter.HandleFunc("/swagger/", docsHandler.ServeUI).Methods(http.MethodGet)
	baseRouter.HandleFunc("/swagger", docsHandler.RedirectToUI).Methods(http.MethodGet)

	api := baseRouter.PathPrefix("/api/v1").Subrouter()

	// ProgressHandlers
	api.HandleFunc("/me/progress", progressHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/me/progress", progressHandler.ListMyProgress).Methods(http.MethodGet)

	// ProblemHandlers
	api.HandleFunc("/problems", problemHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/problems", problemHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/problems/{slug}", problemHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/problems/{slug}", problemHandler.Delete).Methods(http.MethodDelete)

	return router
}
