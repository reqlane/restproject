package handlers

import (
	"encoding/json"
	"net/http"
	"restproject/internal/api/models"
	"restproject/internal/api/services"
	"strconv"
	"time"
)

type execsHandler struct {
	service *services.ExecsService
}

func NewExecsHandler(service *services.ExecsService) *execsHandler {
	return &execsHandler{service: service}
}

// GET /execs/{id}
func (h *execsHandler) GetSingleExecHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid exec id", http.StatusBadRequest)
		return
	}

	exec, err := h.service.GetByID(id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exec)
}

// GET /execs/
func (h *execsHandler) GetExecsHandler(w http.ResponseWriter, r *http.Request) {
	criteria := models.Criteria{
		Filters:  map[string]string{},
		Sortings: r.URL.Query()["sortby"],
	}
	criteria.AddFilters(r.URL.Query(), models.ExecFieldNames)

	execs, err := h.service.GetAllByCriteria(criteria)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(execs),
		Data:   execs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /execs/
func (h *execsHandler) PostExecsHandler(w http.ResponseWriter, r *http.Request) {
	var newExecs []models.Exec
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&newExecs)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	addedExecs, err := h.service.SaveAll(newExecs)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}
	json.NewEncoder(w).Encode(response)
}

// PATCH /execs/{id}
func (h *execsHandler) PatchSingleExecHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid exec id", http.StatusBadRequest)
		return
	}

	var update map[string]any
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedExec, err := h.service.Update(id, update)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedExec)
}

// PATCH /execs/
func (h *execsHandler) PatchExecsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedExecs, err := h.service.UpdateAll(updates)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedExecs)
}

// DELETE /execs/{id}
func (h *execsHandler) DeleteSingleExecHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid exec id", http.StatusBadRequest)
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

// POST /execs/login
func (h *execsHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.ExecCredentials
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&credentials)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tokenString, err := h.service.Login(&credentials)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}
