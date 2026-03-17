package router

import (
	"database/sql"
	"net/http"
	"restproject/internal/api/handlers"
)

type app struct {
	db *sql.DB
}

func NewApp(db *sql.DB) *app {
	return &app{db: db}
}

func (a *app) Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.NewRootHandler(a.db).RootHandler)

	mux.HandleFunc("/teachers/", handlers.NewTeachersHandler(a.db).TeachersHandler)

	mux.HandleFunc("/students/", handlers.NewStudentsHandler(a.db).StudentsHandler)

	mux.HandleFunc("/execs/", handlers.NewExecsHandler(a.db).ExecsHandler)

	return mux
}
