package handlers

import (
	"encoding/json"
	"net/http"
	"restproject/internal/api/models"
	"restproject/internal/api/services"
	"strconv"
)

type teachersHandler struct {
	service *services.TeachersService
}

func NewTeachersHandler(service *services.TeachersService) *teachersHandler {
	return &teachersHandler{service: service}
}

// GET /teachers/{id}
func (h *teachersHandler) GetSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid teacher id", http.StatusBadRequest)
		return
	}

	teacher, err := h.service.GetByID(id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

// GET /teachers/
func (h *teachersHandler) GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	criteria := models.Criteria{
		Filters:  map[string]string{},
		Sortings: r.URL.Query()["sortby"],
	}
	criteria.AddFilters(r.URL.Query(), models.TeacherFieldNames)

	teachers, err := h.service.GetAllByCriteria(criteria)
	if err != nil {
		handleServiceError(w, err)
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

// POST /teachers/
func (h *teachersHandler) PostTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []models.Teacher
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&newTeachers)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	addedTeachers, err := h.service.SaveAll(newTeachers)
	if err != nil {
		handleServiceError(w, err)
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
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid teacher id", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.Replace(id, &updatedTeacher)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// PATCH /teachers/{id}
func (h *teachersHandler) PatchSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid teacher id", http.StatusBadRequest)
		return
	}

	var update map[string]any
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedTeacher, err := h.service.Update(id, update)
	if err != nil {
		handleServiceError(w, err)
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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedTeachers, err := h.service.UpdateAll(updates)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeachers)
}

// DELETE /teachers/{id}
func (h *teachersHandler) DeleteSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid teacher id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		handleServiceError(w, err)
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

// DELETE /teachers/
func (h *teachersHandler) DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	deletedIds, err := h.service.DeleteAll(ids)
	if err != nil {
		handleServiceError(w, err)
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
