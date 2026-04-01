package services

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"restproject/internal/api/models"
	"restproject/internal/apperrors"
	"strconv"
	"strings"

	"github.com/go-mail/mail/v2"
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

func encodePassword(password string) (string, error) {
	if len(password) < 8 {
		return "", apperrors.NewError(apperrors.ErrValidation, errors.New("password must be at least 8 characters"))
	}

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := hash(password, salt)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedPassword := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedPassword, nil
}

func verifyPassword(password string, dbExec *models.Exec) error {
	parts := strings.Split(dbExec.Password, ".")
	if len(parts) != 2 {
		return errors.New("invalid encoded hash format")
	}

	saltBase64 := parts[0]
	dbHashBase64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return err
	}
	dbHash, err := base64.StdEncoding.DecodeString(dbHashBase64)
	if err != nil {
		return err
	}

	requestHash := hash(password, salt)

	if len(requestHash) != len(dbHash) || subtle.ConstantTimeCompare(requestHash, dbHash) != 1 {
		return apperrors.NewError(apperrors.ErrInvalidCredentials, errors.New("invalid credentials"))
	}
	return nil
}

func hash(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

func generatePasswordResetToken() (string, string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", "", err
	}
	token := hex.EncodeToString(tokenBytes)

	hashedTokenBytes := sha256.Sum256(tokenBytes)
	hashedToken := hex.EncodeToString(hashedTokenBytes[:])
	return token, hashedToken, nil
}

func sendEmail(from, to, subject, message string) error {
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	d := mail.NewDialer("localhost", 1025, "", "")
	return d.DialAndSend(m)
}
