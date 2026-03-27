package handlers

import (
	"net/http"
	"restproject/internal/api/services"
)

type execsHandler struct {
	service *services.ExecsService
}

func NewExecsHandler(service *services.ExecsService) *execsHandler {
	return &execsHandler{service: service}
}

func (h *execsHandler) ExecsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Execs Route"))
	case http.MethodPost:
		w.Write([]byte("Hello POST Method on Execs Route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Execs Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Execs Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Execs Route"))
	}
}
