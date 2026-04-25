package handlers

import (
	"github.com/WebCraftersGH/Education-service/internal/contracts"
)

type ProblemContentHandler struct {
	usecase contracts.ProblemContentSVC
}

func NewProblemContentHandler(usecase contracts.ProblemContentSVC) *ProblemContentHandler {
	return &ProblemContentHandler{usecase: usecase}
}
