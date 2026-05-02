package main

import (
	"net/http"
	"os"

	"github.com/WebCraftersGH/Education-service/internal/config"
	"github.com/WebCraftersGH/Education-service/internal/database"
	repo "github.com/WebCraftersGH/Education-service/internal/repository/problem_catalog"
	uc "github.com/WebCraftersGH/Education-service/internal/usecase/problem_catalog"
	"github.com/WebCraftersGH/Education-service/pkg/logging"

	progressRepo "github.com/WebCraftersGH/Education-service/internal/repository/course_progress"
	progressSVC "github.com/WebCraftersGH/Education-service/internal/usecase/course_progress"

	"github.com/WebCraftersGH/Education-service/internal/authclient"
	"github.com/WebCraftersGH/Education-service/internal/middleware"

	transporthttp "github.com/WebCraftersGH/Education-service/internal/transport/http"
	docsHandlers "github.com/WebCraftersGH/Education-service/internal/transport/http/docs"
	handlers "github.com/WebCraftersGH/Education-service/internal/transport/http/handlers"
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
		"app_env":     cfg.AppEnv,
		"http_port":   cfg.HTTPPort,
		"db_host":     cfg.DBHost,
		"db_port":     cfg.DBPort,
		"db_name":     cfg.DBName,
		"db_password": cfg.DBPass,
	}).Info("config loaded")

	db, err := database.NewPostgres(cfg)
	if err != nil {
		logger.WithError(err).Fatal("init database")
	}

	problemRepo := repo.NewRepository(db, logger)
	problemContentRepo := repo.NewRepositoryProblemContent(db, logger)
	problemUC := uc.NewProblemUseCase(problemRepo)
	problemContentUC := uc.NewProblemContentUseCase(problemContentRepo)

	pgRepo := progressRepo.NewProgressRepo(db, logger)
	pgSVC := progressSVC.NewCourseProgress(pgRepo)

	authCl := authclient.New(cfg.AuthServiceURL, logger)

	problemHandler := handlers.NewProblemHandler(problemUC)
	problemContentHandler := handlers.NewProblemContentHandler(problemContentUC)
	progressHandler := handlers.NewProgressHandler(pgSVC, logger)
	docsHandler := docsHandlers.NewDocsHandler()
	healthHandler := handlers.NewHealthHandler()

	router := transporthttp.NewRouter(
		progressHandler,
		problemHandler,
		problemContentHandler,
		healthHandler,
		docsHandler,
		authCl,
		cfg.DEBUG_MODE,
	)

	var handler http.Handler = router

	if cfg.DEBUG_MODE {
		handler = middleware.CORSMiddleware(handler)
	}

	logger.WithField("address", cfg.HTTPAddress()).Info("http server started")
	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.WithError(err).Fatal("http server error")
		os.Exit(1)
	}
}
