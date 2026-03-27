package handlers

import (
	"database/sql"
	"net/http"
)

type rootHandler struct {
	db *sql.DB
}

func NewRootHandler(db *sql.DB) *rootHandler {
	return &rootHandler{db: db}
}

func (h *rootHandler) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to School API"))
}
