package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"restproject/internal/api/models"
	"restproject/internal/apperrors"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

func checkBlankFields(value any) error {
	val := reflect.ValueOf(value)
	for _, field := range val.Fields() {
		if field.Kind() == reflect.String && field.String() == "" {
			return apperrors.NewError(apperrors.ErrValidation, errors.New("all fields are required"))
		}
	}
	return nil
}

func extractID(update map[string]any) (int, error) {
	idRaw, exists := update["id"]
	if !exists {
		return 0, apperrors.NewError(apperrors.ErrMissingID, errors.New("missing id in request body"))
	}

	switch v := idRaw.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, apperrors.NewError(apperrors.ErrInvalidID, fmt.Errorf("invalid id format: '%s' is not a number", v))
		}
		return id, nil
	default:
		return 0, apperrors.NewError(apperrors.ErrInvalidID, errors.New("invalid id format: must be a number"))
	}
}

func applyUpdates(model models.ModelWithID, update map[string]any) error {
	modelVal := reflect.ValueOf(model).Elem()
	modelType := modelVal.Type()

	for k, v := range update {
		if k == "id" {
			continue
		}
		for i := 0; i < modelVal.NumField(); i++ {
			typeField := modelType.Field(i)
			valField := modelVal.Field(i)
			jsonName := strings.Split(typeField.Tag.Get("json"), ",")[0]
			if jsonName == k {
				if valField.CanSet() {
					value := reflect.ValueOf(v)
					if value.Type().ConvertibleTo(typeField.Type) {
						valField.Set(value.Convert(typeField.Type))
					} else {
						return apperrors.NewError(apperrors.ErrInvalidField, fmt.Errorf("invalid type for field '%s' on id %d", k, model.GetID()))
					}
				}
				break
			}
		}
	}
	return nil
}

func encodePassword(model models.ModelWithPassword) error {
	if len(model.GetPassword()) < 8 {
		return apperrors.NewError(apperrors.ErrValidation, errors.New("password must be at least 8 characters"))
	}

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	hash := argon2.IDKey([]byte(model.GetPassword()), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	model.SetPassword(encodedHash)
	return nil
}
