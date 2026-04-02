package handlers

import (
	"log"
	"net/http"
	"restproject/internal/api/models"
	"restproject/internal/apperrors"
	"strconv"
)

func handleServiceError(w http.ResponseWriter, err error) {
	httpErr := apperrors.FromError(err)
	if httpErr.Status == http.StatusInternalServerError {
		log.Println(err)
	}
	http.Error(w, httpErr.Message, httpErr.Status)
}

func paginationFrom(r *http.Request) *models.Pagination {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	return &models.Pagination{Page: page, Limit: limit}
}
