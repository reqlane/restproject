package handlers

import (
	"log"
	"net/http"
	"restproject/internal/apperrors"
)

func handleServiceError(w http.ResponseWriter, err error) {
	httpErr := apperrors.FromError(err)
	if httpErr.Status == http.StatusInternalServerError {
		log.Println(err)
	}
	http.Error(w, httpErr.Message, httpErr.Status)
}
