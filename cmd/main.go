package main

import (
	"net/http"

	"github.com/WebCraftersGH/Education-service/internal/config"
	controller "github.com/WebCraftersGH/Education-service/internal/controller/problem_catalog"
	"github.com/WebCraftersGH/Education-service/internal/database"
	appmiddleware "github.com/WebCraftersGH/Education-service/internal/middleware"
	repo "github.com/WebCraftersGH/Education-service/internal/repository/problem_catalog"
	uc "github.com/WebCraftersGH/Education-service/internal/usecase/problem_catalog"
	"github.com/WebCraftersGH/Education-service/pkg/logging"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		panic(err)
	}

	logger, closer, err := logging.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

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

	problemRepo := repo.NewRepository(db, logger)
	problemContentRepo := repo.NewRepositoryProblemContent(db, logger)

	problemUC := uc.NewProblemUseCase(problemRepo)
	problemContentUC := uc.NewProblemContentUseCase(problemContentRepo)

	problemCatalogCTRL := controller.NewProblemCatalogController(
		problemUC,
		problemContentUC,
		logger,
	)

	r := chi.NewRouter()

	r.Use(appmiddleware.GenerateRequestID)

	problemCatalogCTRL.RegisterRoutes(r)

	logger.WithField("address", cfg.HTTPAddress()).Info("http server started")
	if err := http.ListenAndServe(cfg.HTTPAddress(), r); err != nil {
		logger.WithError(err).Fatal("http server stopped")
	}
}
