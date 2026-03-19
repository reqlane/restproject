package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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

// GET /teachers/
func (h *teachersHandler) GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var query strings.Builder
	query.WriteString(`SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1`)
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

	err = rows.Err()
	if err != nil {
		http.Error(w, "Error during row iteration", http.StatusInternalServerError)
		return
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
}

// GET /teachers/{id}
func (h *teachersHandler) GetSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var teacher models.Teacher
	query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id=?`
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
	addedSort := false
	if len(sortParams) > 0 {
		for _, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			if !addedSort {
				query.WriteString(" ORDER BY")
				addedSort = true
			} else {
				query.WriteString(",")
			}
			query.WriteString(" " + field + " " + order)
		}
	}
}

func isValidSortOrder(order string) bool {
	orderLowerCase := strings.ToLower(order)
	return orderLowerCase == "asc" || orderLowerCase == "desc"
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
func (h *teachersHandler) PostTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)`)
	if err != nil {
		http.Error(w, "Error in preparing SQL query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
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

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
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
func (h *teachersHandler) PutSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

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
	query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id=?`
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
	query = `UPDATE teachers SET first_name=?, last_name=?, email=?, class=?, subject=? WHERE id=?`
	_, err = h.db.Exec(query, updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

// PATCH /teachers/
func (h *teachersHandler) PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	updateStmt, err := tx.Prepare(`UPDATE teachers SET first_name=?, last_name=?, email=?, class=?, subject=? WHERE id=?`)
	if err != nil {
		http.Error(w, "Error preparing update statement", http.StatusInternalServerError)
		return
	}
	defer updateStmt.Close()

	updatedTeachers := make([]models.Teacher, 0, len(updates))

	for _, update := range updates {
		idRaw, exists := update["id"]
		if !exists {
			http.Error(w, "Teacher ID required", http.StatusBadRequest)
			return
		}

		var id int
		switch v := idRaw.(type) {
		case float64:
			id = int(v)
		case int:
			id = v
		case string:
			id, err = strconv.Atoi(v)
			if err != nil {
				http.Error(w, "Error converting ID to int", http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var dbTeacher models.Teacher
		query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?`
		err = tx.QueryRow(query, id).Scan(&dbTeacher.ID, &dbTeacher.FirstName, &dbTeacher.LastName, &dbTeacher.Email, &dbTeacher.Class, &dbTeacher.Subject)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Teacher not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Error retrieving teacher", http.StatusInternalServerError)
			return
		}

		// Apply updates using reflection
		teacherVal := reflect.ValueOf(&dbTeacher).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				typeField := teacherType.Field(i)
				valField := teacherVal.Field(i)
				jsonName := strings.Split(typeField.Tag.Get("json"), ",")[0]
				if jsonName == k {
					if valField.CanSet() {
						value := reflect.ValueOf(v)
						if value.Type().ConvertibleTo(typeField.Type) {
							valField.Set(value.Convert(typeField.Type))
						} else {
							http.Error(w, fmt.Sprintf("Invalid type for field %s", k), http.StatusBadRequest)
							return
						}
					}
					break
				}
			}
		}

		_, err = updateStmt.Exec(dbTeacher.FirstName, dbTeacher.LastName, dbTeacher.Email, dbTeacher.Class, dbTeacher.Subject, dbTeacher.ID)
		if err != nil {
			http.Error(w, "Error updating teacher", http.StatusInternalServerError)
			return
		}

		updatedTeachers = append(updatedTeachers, dbTeacher)
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeachers)
}

// PATCH /teachers/{id}
func (h *teachersHandler) PatchSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	var update map[string]any
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var dbTeacher models.Teacher
	query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id=?`
	err = h.db.QueryRow(query, id).Scan(&dbTeacher.ID, &dbTeacher.FirstName, &dbTeacher.LastName, &dbTeacher.Email, &dbTeacher.Class, &dbTeacher.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error retrieving teacher", http.StatusInternalServerError)
		return
	}

	// Apply updates using reflection
	teacherVal := reflect.ValueOf(&dbTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range update {
		if k == "id" {
			continue
		}
		for i := 0; i < teacherVal.NumField(); i++ {
			typeField := teacherType.Field(i)
			valField := teacherVal.Field(i)
			jsonName := strings.Split(typeField.Tag.Get("json"), ",")[0]
			if jsonName == k {
				if valField.CanSet() {
					value := reflect.ValueOf(v)
					if value.Type().ConvertibleTo(typeField.Type) {
						valField.Set(value.Convert(typeField.Type))
					} else {
						http.Error(w, fmt.Sprintf("Invalid type for field %s", k), http.StatusBadRequest)
						return
					}
				}
				break
			}
		}
	}

	query = `UPDATE teachers SET first_name=?, last_name=?, email=?, class=?, subject=? WHERE id=?`
	_, err = h.db.Exec(query, dbTeacher.FirstName, dbTeacher.LastName, dbTeacher.Email, dbTeacher.Class, dbTeacher.Subject, dbTeacher.ID)
	if err != nil {
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbTeacher)
}

// DELETE /teachers/
func (h *teachersHandler) DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {

}

// DELETE /teachers/{id}
func (h *teachersHandler) DeleteSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM teachers WHERE id=?`
	result, err := h.db.Exec(query, id)
	if err != nil {
		http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error retrieving delete result", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	// No response body
	// w.WriteHeader(http.StatusNoContent)

	// Response Body
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Teacher successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}
