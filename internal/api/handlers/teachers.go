package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"restproject/internal/models"
	"strconv"
	"strings"
)

type teachersHandler struct {
	db *sql.DB
}

func NewTeachersHandler(db *sql.DB) *teachersHandler {
	return &teachersHandler{db: db}
}

func (h *teachersHandler) TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getTeachersHandler(w, r)
	case http.MethodPost:
		h.postTeachersHandler(w, r)
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Teachers Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Teachers Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Teachers Route"))
	}
}

func (h *teachersHandler) getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1 = 1`
		var args []any
		if firstName != "" {
			query += " AND first_name = ?"
			args = append(args, firstName)
		}
		if lastName != "" {
			query += " AND last_name = ?"
			args = append(args, lastName)
		}

		rows, err := h.db.Query(query, args...)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		teacherList := make([]models.Teacher, 0)
		for rows.Next() {
			var teacher models.Teacher
			err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
			if err != nil {
				http.Error(w, "Error scanning database results", http.StatusInternalServerError)
				return
			}
			teacherList = append(teacherList, teacher)
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		// Handle Path parameter
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var teacher models.Teacher
		query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?`
		err = h.db.QueryRow(query, id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teacher)
	}
}

func (h *teachersHandler) postTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	statement, err := h.db.Prepare(`INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)`)
	if err != nil {
		http.Error(w, "Error in preparing SQL query", http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := statement.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, "Error inserting data into database", http.StatusInternalServerError)
			return
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting last inserted ID", http.StatusInternalServerError)
			return
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
}
