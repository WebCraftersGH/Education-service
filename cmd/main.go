package main

import (
	"net/http"

	controller "github.com/WebCraftersGH/Education-service/internal/controller/problem_catalog"
	repo "github.com/WebCraftersGH/Education-service/internal/repository/problem_catalog"
	uc "github.com/WebCraftersGH/Education-service/internal/usecase/problem_catalog"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func main() {
	var db *gorm.DB

	problemRepo := repo.NewRepository(db)
	problemContentRepo := repo.NewRepository(db)

	problemUC := uc.NewProblemUseCase(problemRepo)
	problemContentUC := uc.NewProblemContentUseCase(problemContentRepo)

	problemCatalogCTRL := controller.NewProblemCatalogController(problemUC, problemContentUC)

	r := chi.NewRouter()
	problemCatalogCTRL.RegisterRoutes(r)

	_ = http.ListenAndServe(":8080", r)
}
