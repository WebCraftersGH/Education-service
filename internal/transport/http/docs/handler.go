package docs

import "net/http"

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

func (h *DocsHandler) ServeSpec(w http.ResponseWriter, r *http.Request)    {}
func (h *DocsHandler) ServeUI(w http.ResponseWriter, r *http.Request)      {}
func (h *DocsHandler) RedirectToUI(w http.ResponseWriter, r *http.Request) {}
