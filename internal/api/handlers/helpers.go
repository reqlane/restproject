package handlers

import (
	"log"
	"net/http"
	"restproject/internal/api/models"
	"restproject/internal/apperrors"
)

func handleServiceError(w http.ResponseWriter, err error) {
	httpErr := apperrors.FromError(err)
	if httpErr.Status == http.StatusInternalServerError {
		log.Println(err)
	}
	http.Error(w, httpErr.Message, httpErr.Status)
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
