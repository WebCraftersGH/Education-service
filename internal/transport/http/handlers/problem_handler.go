package handlers

import (
	"github.com/WebCraftersGH/Education-service/internal/contracts"
	"net/http"
)

type ProblemHandler struct {
	usecase contracts.ProblemSVC
}

func NewProblemHandler(usecase contracts.ProblemSVC) *ProblemHandler {
	return &ProblemHandler{usecase: usecase}
}

func (h *ProblemHandler) Create(w http.ResponseWriter, r *http.Request) {}
func (h *ProblemHandler) List(w http.ResponseWriter, r *http.Request)   {}
func (h *ProblemHandler) Update(w http.ResponseWriter, r *http.Request) {}
func (h *ProblemHandler) Delete(w http.ResponseWriter, r *http.Request) {}
