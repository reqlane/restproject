package handlers

import (
	"encoding/json"
	"fmt"
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

func (h *execsHandler) GetExecsHandler(w http.ResponseWriter, r *http.Request) {
	pg := paginationFrom(r)
	criteria := &models.Criteria{
		Filters:  map[string]string{},
		Sortings: r.URL.Query()["sortby"],
	}
	criteria.AddFilters(r.URL.Query(), models.ExecFieldNames)

	execs, totalCount, err := h.service.GetAllByCriteria(criteria, pg)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response := struct {
		Status   string                `json:"status"`
		Count    int                   `json:"count"`
		Page     int                   `json:"page"`
		PageSize int                   `json:"page_size"`
		Data     []models.ExecResponse `json:"data"`
	}{
		Status:   "success",
		Count:    totalCount,
		Page:     pg.Page,
		PageSize: pg.Limit,
		Data:     execs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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
		Status string                `json:"status"`
		Count  int                   `json:"count"`
		Data   []models.ExecResponse `json:"data"`
	}{
		Status: "success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}
	json.NewEncoder(w).Encode(response)
}

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

func (h *execsHandler) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid exec id", http.StatusBadRequest)
		return
	}

	var req models.UpdatePasswordRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tokenString, err := h.service.UpdatePassword(id, &req)
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
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *execsHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ExecCredentials
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tokenString, err := h.service.Login(&req)
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
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *execsHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *execsHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.ForgotPassword(req.Email)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	fmt.Fprintf(w, "password reset link sent to %s", req.Email)
}

func (h *execsHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("resetcode")

	var req models.ResetPasswordRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.Token = token

	err = h.service.ResetPassword(&req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	fmt.Fprint(w, "password reset successfully")
}
