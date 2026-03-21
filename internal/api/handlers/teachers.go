package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"restproject/internal/models"
	"restproject/internal/repository/repositories"
	"strconv"
	"strings"
)

type teachersHandler struct {
	repo *repositories.TeacherRepository
}

func NewTeachersHandler(db *sql.DB) *teachersHandler {
	return &teachersHandler{repo: repositories.NewTeacherRepository(db)}
}

// GET /teachers/
func (h *teachersHandler) GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	criteria := models.TeacherCriteria{
		Filters:  map[string]string{},
		Sortings: r.URL.Query()["sortby"],
	}
	addFiltersCriteria(r, &criteria)

	teachers, err := h.repo.GetAllByCriteria(criteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachers),
		Data:   teachers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func addFiltersCriteria(r *http.Request, criteria *models.TeacherCriteria) {
	fieldNames := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range fieldNames {
		value := r.URL.Query().Get(param)
		if value != "" {
			criteria.Filters[dbField] = value
		}
	}
}

// GET /teachers/{id}
func (h *teachersHandler) GetSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	teacher, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

// POST /teachers/
func (h *teachersHandler) PostTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	addedTeachers, err := h.repo.SaveAll(newTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	dbTeacher, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedTeacher.ID = dbTeacher.ID
	result, err := h.repo.Update(&updatedTeacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// PATCH /teachers/
func (h *teachersHandler) PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

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

		dbTeacher, err := h.repo.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		teacherVal := reflect.ValueOf(dbTeacher).Elem()
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

		updatedTeachers = append(updatedTeachers, *dbTeacher)
	}

	updatedTeachers, err = h.repo.UpdateAll(updatedTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	dbTeacher, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	teacherVal := reflect.ValueOf(dbTeacher).Elem()
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

	updatedTeacher, err := h.repo.Update(dbTeacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

// DELETE /teachers/
func (h *teachersHandler) DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deletedIds, err := h.repo.DeleteAll(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(deletedIds) < 1 {
		http.Error(w, "IDs do not exist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status     string `json:"status"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{
		Status:     "success",
		DeletedIDs: deletedIds,
	}
	json.NewEncoder(w).Encode(response)
}

// DELETE /teachers/{id}
func (h *teachersHandler) DeleteSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "success",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}
