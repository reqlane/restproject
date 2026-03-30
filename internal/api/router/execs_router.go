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

	mux.HandleFunc("GET /execs", execsHandler.GetExecsHandler)
	mux.HandleFunc("POST /execs", execsHandler.PostExecsHandler)
	mux.HandleFunc("PATCH /execs", execsHandler.PatchExecsHandler)

	mux.HandleFunc("GET /execs/{id}", execsHandler.GetSingleExecHandler)
	mux.HandleFunc("PATCH /execs/{id}", execsHandler.PatchSingleExecHandler)
	mux.HandleFunc("DELETE /execs/{id}", execsHandler.DeleteSingleExecHandler)
	mux.HandleFunc("POST /execs/{id}/updatepassword", execsHandler.GetExecsHandler)

	mux.HandleFunc("POST /execs/login", execsHandler.LoginHandler)
	mux.HandleFunc("POST /execs/logout", execsHandler.LogoutHandler)
	mux.HandleFunc("POST /execs/forgotpassword", execsHandler.GetExecsHandler)
	mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", execsHandler.GetExecsHandler)
}
