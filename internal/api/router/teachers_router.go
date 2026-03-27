package router

import (
	"net/http"
	"restproject/internal/api/handlers"
	"restproject/internal/api/repositories"
	"restproject/internal/api/services"
)

func (a *app) teachersRouter(mux *http.ServeMux) {
	teachersRepo := repositories.NewTeacherRepository(a.db)
	teachersService := services.NewTeachersService(teachersRepo)
	teachersHandler := handlers.NewTeachersHandler(teachersService)

	mux.HandleFunc("GET /teachers", teachersHandler.GetTeachersHandler)
	mux.HandleFunc("POST /teachers", teachersHandler.PostTeachersHandler)
	mux.HandleFunc("PATCH /teachers", teachersHandler.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers", teachersHandler.DeleteTeachersHandler)

	mux.HandleFunc("GET /teachers/{id}", teachersHandler.GetSingleTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", teachersHandler.PutSingleTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", teachersHandler.PatchSingleTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", teachersHandler.DeleteSingleTeacherHandler)

	mux.HandleFunc("GET /teachers/{id}/students", teachersHandler.GetStudentsByTeacherID)
	mux.HandleFunc("GET /teachers/{id}/students/count", teachersHandler.GetStudentsCountByTeacherID)
}
