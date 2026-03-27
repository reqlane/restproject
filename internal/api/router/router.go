package router

import (
	"database/sql"
	"net/http"
	"restproject/internal/api/handlers"
	"restproject/internal/api/repositories"
	"restproject/internal/api/services"
)

type app struct {
	db *sql.DB
}

func NewApp(db *sql.DB) *app {
	return &app{db: db}
}

func (a *app) Router() *http.ServeMux {
	mux := http.NewServeMux()

	// Root
	rootHandler := handlers.NewRootHandler(a.db)
	mux.HandleFunc("GET /", rootHandler.RootHandler)

	// Teachers
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
	mux.HandleFunc("GET /teachers/{id}/studentscount", teachersHandler.GetStudentsCountByTeacherID)

	// Students
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

	//Execs
	execsHandler := handlers.NewExecsHandler(a.db)
	mux.HandleFunc("GET /execs/", execsHandler.ExecsHandler)

	return mux
}
