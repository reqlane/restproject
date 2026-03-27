package router

import (
	"net/http"
	"restproject/internal/api/handlers"
	"restproject/internal/api/repositories"
	"restproject/internal/api/services"
)

func (a *app) execsRouter(mux *http.ServeMux) {
	execsRepository := repositories.NewExecsRepository(a.db)
	execsService := services.NewExecsService(execsRepository)
	execsHandler := handlers.NewExecsHandler(execsService)

	mux.HandleFunc("GET /execs", execsHandler.ExecsHandler)
	mux.HandleFunc("POST /execs", execsHandler.ExecsHandler)
	mux.HandleFunc("PATCH /execs", execsHandler.ExecsHandler)

	mux.HandleFunc("GET /execs/{id}", execsHandler.ExecsHandler)
	mux.HandleFunc("PATCH /execs/{id}", execsHandler.ExecsHandler)
	mux.HandleFunc("DELETE /execs/{id}", execsHandler.ExecsHandler)
	mux.HandleFunc("POST /execs/{id}/updatepassword", execsHandler.ExecsHandler)

	mux.HandleFunc("POST /execs/login", execsHandler.ExecsHandler)
	mux.HandleFunc("POST /execs/logout", execsHandler.ExecsHandler)
	mux.HandleFunc("POST /execs/forgotpassword", execsHandler.ExecsHandler)
	mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", execsHandler.ExecsHandler)
}
