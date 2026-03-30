package routers

import (
	"net/http"
	"restproject/internal/api/handlers"
	"restproject/internal/api/repositories"
	"restproject/internal/api/services"
)

func (a *app) studentsRouter(mux *http.ServeMux) {
	studentsRepo := repositories.NewStudentsRepository(a.db)
	studentsService := services.NewStudentsService(studentsRepo)
	studentsHandler := handlers.NewStudentsHandler(studentsService)

	mux.HandleFunc("GET /students", studentsHandler.GetStudentsHandler)
	mux.HandleFunc("POST /students", studentsHandler.PostStudentsHandler)
	mux.HandleFunc("PATCH /students", studentsHandler.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students", studentsHandler.DeleteStudentsHandler)

	mux.HandleFunc("GET /students/{id}", studentsHandler.GetSingleStudentHandler)
	mux.HandleFunc("PUT /students/{id}", studentsHandler.PutSingleStudentHandler)
	mux.HandleFunc("PATCH /students/{id}", studentsHandler.PatchSingleStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", studentsHandler.DeleteSingleStudentHandler)
}
