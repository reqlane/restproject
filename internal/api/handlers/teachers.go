package handlers

import (
	"encoding/json"
	"net/http"
	"restproject/internal/api/models"
	"restproject/internal/api/services"
	"restproject/internal/auth"
	"strconv"
)

type teachersHandler struct {
	service *services.TeachersService
}

func NewTeachersHandler(service *services.TeachersService) *teachersHandler {
	return &teachersHandler{service: service}
}

func (h *teachersHandler) GetSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.AuthorizeUser(r.Context().Value(auth.ContextKeyRole).(string), auth.Admin, auth.Manager, auth.Exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

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

func (h *teachersHandler) GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.AuthorizeUser(r.Context().Value(auth.ContextKeyRole).(string), auth.Admin, auth.Manager, auth.Exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	pg := paginationFrom(r)
	criteria := models.Criteria{
		Filters:  map[string]string{},
		Sortings: r.URL.Query()["sortby"],
	}
	criteria.AddFilters(r.URL.Query(), models.TeacherFieldNames)

	teachers, totalCount, err := h.service.GetAllByCriteria(criteria, pg)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response := struct {
		Status   string           `json:"status"`
		Count    int              `json:"count"`
		Page     int              `json:"page"`
		PageSize int              `json:"page_size"`
		Data     []models.Teacher `json:"data"`
	}{
		Status:   "success",
		Count:    totalCount,
		Page:     pg.Page,
		PageSize: pg.Limit,
		Data:     teachers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *teachersHandler) PostTeachersHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.AuthorizeUser(r.Context().Value(auth.ContextKeyRole).(string), auth.Admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var newTeachers []models.Teacher
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&newTeachers)
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

func (h *teachersHandler) PutSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.AuthorizeUser(r.Context().Value(auth.ContextKeyRole).(string), auth.Admin, auth.Exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

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

func (h *teachersHandler) PatchSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.AuthorizeUser(r.Context().Value(auth.ContextKeyRole).(string), auth.Admin, auth.Exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

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

func (h *teachersHandler) PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.AuthorizeUser(r.Context().Value(auth.ContextKeyRole).(string), auth.Admin, auth.Exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var updates []map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
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

func (h *teachersHandler) DeleteSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.AuthorizeUser(r.Context().Value(auth.ContextKeyRole).(string), auth.Admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

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

func (h *teachersHandler) DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.AuthorizeUser(r.Context().Value(auth.ContextKeyRole).(string), auth.Admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var ids []int
	err = json.NewDecoder(r.Body).Decode(&ids)
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

func (h *teachersHandler) GetStudentsByTeacherID(w http.ResponseWriter, r *http.Request) {
	err := auth.AuthorizeUser(r.Context().Value(auth.ContextKeyRole).(string), auth.Admin, auth.Manager, auth.Exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid teacher id", http.StatusBadRequest)
		return
	}

	students, err := h.service.GetStudentsByTeacherID(id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(students),
		Data:   students,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *teachersHandler) GetStudentsCountByTeacherID(w http.ResponseWriter, r *http.Request) {
	err := auth.AuthorizeUser(r.Context().Value(auth.ContextKeyRole).(string), auth.Admin, auth.Manager, auth.Exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid teacher id", http.StatusBadRequest)
		return
	}

	studentsCount, err := h.service.GetStudentsCountByTeacherID(id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}{
		Status: "success",
		Count:  studentsCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
