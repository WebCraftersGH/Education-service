package handlers

import "net/http"

type ProblemHandler struct{}

func NewProblemHandler() *ProgressHandler {
	return &ProgressHandler{}
}

func (h *ProblemHandler) Create(w http.ResponseWriter, r *http.Request) {}
func (h *ProblemHandler) List(w http.ResponseWriter, r *http.Request)   {}
func (h *ProblemHandler) Update(w http.ResponseWriter, r *http.Request) {}
func (h *ProblemHandler) Delete(w http.ResponseWriter, r *http.Request) {}
