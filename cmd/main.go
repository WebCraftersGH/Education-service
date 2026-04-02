package main

import (
	"net/http"
	"os"

	"github.com/WebCraftersGH/Education-service/internal/config"
	controller "github.com/WebCraftersGH/Education-service/internal/controller/problem_catalog"
	"github.com/WebCraftersGH/Education-service/internal/database"
	appmiddleware "github.com/WebCraftersGH/Education-service/internal/middleware"
	repo "github.com/WebCraftersGH/Education-service/internal/repository/problem_catalog"
	uc "github.com/WebCraftersGH/Education-service/internal/usecase/problem_catalog"
	"github.com/WebCraftersGH/Education-service/pgk/logging"
	"github.com/go-chi/chi/v5"
)

func main() {
	logging.Init(os.Getenv("LOG_LEVEL"))
	logger := logging.GetLogger()

	cfg, err := config.Load(".env")
	if err != nil {
		logger.WithError(err).Fatal("load config")
	}

	logging.Init(cfg.LogLevel)
	logger = logging.GetLogger()
	logger.WithFields(map[string]any{
		"app_env":   cfg.AppEnv,
		"http_port": cfg.HTTPPort,
		"db_host":   cfg.DBHost,
		"db_port":   cfg.DBPort,
		"db_name":   cfg.DBName,
	}).Info("config loaded")

	db, err := database.NewPostgres(cfg)
	if err != nil {
		logger.WithError(err).Fatal("init database")
	}

	problemRepo := repo.NewRepository(db)
	problemContentRepo := repo.NewRepositoryProblemContent(db)

	problemUC := uc.NewProblemUseCase(problemRepo)
	problemContentUC := uc.NewProblemContentUseCase(problemContentRepo)

	problemCatalogCTRL := controller.NewProblemCatalogController(problemUC, problemContentUC)

	r := chi.NewRouter()
	r.Use(appmiddleware.RequestLogger)
	problemCatalogCTRL.RegisterRoutes(r)

	logger.WithField("address", cfg.HTTPAddress()).Info("http server started")
	if err := http.ListenAndServe(cfg.HTTPAddress(), r); err != nil {
		logger.WithError(err).Fatal("http server stopped")
	}
}
