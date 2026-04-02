package handlers

import (
	"encoding/json"
	"net/http"
	"restproject/internal/api/models"
	"restproject/internal/api/services"
	"strconv"
)

type studentsHandler struct {
	service *services.StudentsService
}

func NewStudentsHandler(service *services.StudentsService) *studentsHandler {
	return &studentsHandler{service: service}
}

func (h *studentsHandler) GetSingleStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid student id", http.StatusBadRequest)
		return
	}

	student, err := h.service.GetByID(id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func (h *studentsHandler) GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	pg := paginationFrom(r)
	criteria := &models.Criteria{
		Filters:  map[string]string{},
		Sortings: r.URL.Query()["sortby"],
	}
	criteria.AddFilters(r.URL.Query(), models.StudentFieldNames)

	students, totalCount, err := h.service.GetAllByCriteria(criteria, pg)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response := struct {
		Status   string           `json:"status"`
		Count    int              `json:"count"`
		Page     int              `json:"page"`
		PageSize int              `json:"page_size"`
		Data     []models.Student `json:"data"`
	}{
		Status:   "success",
		Count:    totalCount,
		Page:     pg.Page,
		PageSize: pg.Limit,
		Data:     students,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *studentsHandler) PostStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var newStudents []models.Student
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&newStudents)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	addedStudents, err := h.service.SaveAll(newStudents)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addedStudents),
		Data:   addedStudents,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *studentsHandler) PutSingleStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid student id", http.StatusBadRequest)
		return
	}

	var updatedStudent models.Student
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&updatedStudent)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.Replace(id, &updatedStudent)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *studentsHandler) PatchSingleStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid student id", http.StatusBadRequest)
		return
	}

	var update map[string]any
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedStudent, err := h.service.Update(id, update)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudent)
}

func (h *studentsHandler) PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedStudents, err := h.service.UpdateAll(updates)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudents)
}

func (h *studentsHandler) DeleteSingleStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid student id", http.StatusBadRequest)
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

func (h *studentsHandler) DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {
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
