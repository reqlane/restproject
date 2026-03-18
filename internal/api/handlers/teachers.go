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
		h.updateTeachersHandler(w, r)
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Teachers Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Teachers Route"))
	}
}

// GET /teachers/
func (h *teachersHandler) getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		var query strings.Builder
		query.WriteString(`SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1 = 1`)
		var args []any

		args = addFilters(r, &query, args)

		addSorting(r, &query)

		rows, err := h.db.Query(query.String(), args...)
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

func addFilters(r *http.Request, query *strings.Builder, args []any) []any {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query.WriteString(" AND " + dbField + " = ?")
			args = append(args, value)
		}
	}
	return args
}

func addSorting(r *http.Request, query *strings.Builder) {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		query.WriteString(" ORDER BY")
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query.WriteString(",")
			}
			query.WriteString(" " + field + " " + order)
		}
	}
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

// POST /teachers/
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

// PUT /teachers/{id}
func (h *teachersHandler) updateTeachersHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var existingTeacher models.Teacher
	query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?`
	err = h.db.QueryRow(query, id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}

	updatedTeacher.ID = existingTeacher.ID
	query = `UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?`
	_, err = h.db.Exec(query, &updatedTeacher.FirstName, &updatedTeacher.LastName, &updatedTeacher.Email, &updatedTeacher.Class, &updatedTeacher.Subject, &updatedTeacher.ID)
	if err != nil {
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}
