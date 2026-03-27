package router

import (
	"net/http"
	"restproject/internal/api/handlers"
)

func (a *app) execsRouter(mux *http.ServeMux) {
	execsHandler := handlers.NewExecsHandler(a.db)
	mux.HandleFunc("GET /execs/", execsHandler.ExecsHandler)
}
