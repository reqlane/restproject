package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"restproject/internal/api/models"
	"restproject/internal/api/repositories"
	"restproject/internal/apperrors"
	"restproject/internal/auth"
	"strconv"
	"time"
)

type ExecsService struct {
	repo *repositories.ExecsRepository
}

func NewExecsService(repo *repositories.ExecsRepository) *ExecsService {
	return &ExecsService{repo: repo}
}

func (s *ExecsService) GetByID(id int) (*models.ExecResponse, error) {
	exec, err := s.getByID(id)
	if err != nil {
		return nil, fmt.Errorf("service.GetByID: %w", err)
	}
	return exec.ToResponse(), nil
}

func (s *ExecsService) getByID(id int) (*models.Exec, error) {
	exec, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewError(apperrors.ErrNotFound, fmt.Errorf("exec id %d not found", id))
		}
		return nil, err
	}
	return exec, nil
}

func (s *ExecsService) GetAllByCriteria(criteria *models.Criteria, pg *models.Pagination) ([]models.ExecResponse, int, error) {
	execs, totalCount, err := s.repo.GetAllByCriteria(criteria, pg)
	if err != nil {
		return nil, 0, fmt.Errorf("service.GetAllByCriteria: %w", err)
	}
	return models.Execs(execs).ToResponse(), totalCount, nil
}

func (s *ExecsService) SaveAll(execs []models.Exec) ([]models.ExecResponse, error) {
	for i, exec := range execs {
		if err := checkBlankFields(exec); err != nil {
			return nil, fmt.Errorf("service.SaveAll: %w", err)
		}
		encodedPassword, err := encodePassword(exec.Password)
		if err != nil {
			return nil, fmt.Errorf("service.SaveAll: %w", err)
		}
		execs[i].Password = encodedPassword
	}

	execs, err := s.repo.SaveAll(execs)
	if err != nil {
		return nil, fmt.Errorf("service.SaveAll: %w", err)
	}
	return models.Execs(execs).ToResponse(), nil
}

func (s *ExecsService) Update(id int, update map[string]any) (*models.ExecResponse, error) {
	dbExec, err := s.getByID(id)
	if err != nil {
		return nil, fmt.Errorf("service.Update: %w", err)
	}

	if err = applyUpdates(dbExec, update); err != nil {
		return nil, fmt.Errorf("service.Update: %w", err)
	}

	updatedExec, err := s.repo.Update(dbExec)
	if err != nil {
		return nil, fmt.Errorf("service.Update: %w", err)
	}
	return updatedExec.ToResponse(), nil
}

func (s *ExecsService) UpdateAll(updates []map[string]any) ([]models.ExecResponse, error) {
	updatedExecs := make([]models.Exec, 0, len(updates))

	for _, update := range updates {
		id, err := extractID(update)
		if err != nil {
			return nil, fmt.Errorf("service.UpdateAll: %w", err)
		}

		dbExec, err := s.getByID(id)
		if err != nil {
			return nil, fmt.Errorf("service.UpdateAll: %w", err)
		}

		if err = applyUpdates(dbExec, update); err != nil {
			return nil, fmt.Errorf("service.UpdateAll: %w", err)
		}

		updatedExecs = append(updatedExecs, *dbExec)
	}

	updatedExecs, err := s.repo.UpdateAll(updatedExecs)
	if err != nil {
		return nil, fmt.Errorf("service.UpdateAll: %w", err)
	}
	return models.Execs(updatedExecs).ToResponse(), nil
}

func (s *ExecsService) Delete(id int) error {
	_, err := s.getByID(id)
	if err != nil {
		return fmt.Errorf("service.Delete: %w", err)
	}

	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("service.Delete: %w", err)
	}
	return nil
}

func (s *ExecsService) UpdatePassword(id int, req *models.UpdatePasswordRequest) (string, error) {
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return "", apperrors.NewError(apperrors.ErrValidation, errors.New("current and new passwords are required"))
	}

	exec, err := s.getByID(id)
	if err != nil {
		return "", fmt.Errorf("service.UpdatePasword: %w", err)
	}

	err = verifyPassword(req.CurrentPassword, exec)
	if err != nil {
		return "", fmt.Errorf("service.UpdatePasword: %w", err)
	}

	encodedPassword, err := encodePassword(req.NewPassword)
	if err != nil {
		return "", fmt.Errorf("service.UpdatePasword: %w", err)
	}

	currentTime := time.Now().Format(time.RFC3339)
	exec.PasswordChangedAt = sql.NullString{String: currentTime, Valid: true}
	exec.Password = encodedPassword

	err = s.repo.UpdatePassword(exec)
	if err != nil {
		return "", fmt.Errorf("service.UpdatePasword: %w", err)
	}

	tokenString, err := auth.SignToken(exec.ID, exec.Username, exec.Role)
	if err != nil {
		return "", fmt.Errorf("service.Login: %w", err)
	}
	return tokenString, nil
}

func (s *ExecsService) Login(credentials *models.ExecCredentials) (string, error) {
	if credentials.Username == "" || credentials.Password == "" {
		return "", fmt.Errorf("service.Login: %w", apperrors.NewError(apperrors.ErrValidation, errors.New("username and password are required")))
	}

	exec, err := s.repo.GetByUsername(credentials.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("service.Login: %w", apperrors.NewError(apperrors.ErrInvalidCredentials, errors.New("invalid credentials")))
		}
		return "", fmt.Errorf("service.Login: %w", err)
	}

	if exec.InactiveStatus {
		return "", fmt.Errorf("service.Login: %w", apperrors.NewError(apperrors.ErrInactiveAccount, errors.New("account is inactive")))
	}

	err = verifyPassword(credentials.Password, exec)
	if err != nil {
		return "", fmt.Errorf("service.Login: %w", err)
	}

	tokenString, err := auth.SignToken(exec.ID, exec.Username, exec.Role)
	if err != nil {
		return "", fmt.Errorf("service.Login: %w", err)
	}

	return tokenString, nil
}

func (s *ExecsService) ForgotPassword(email string) error {
	if email == "" {
		return fmt.Errorf("service.ForgotPassword: %w", apperrors.NewError(apperrors.ErrValidation, errors.New("user email is required")))
	}
	exec, err := s.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("service.ForgotPassword: %w", err)
	}

	resetTokenDurationMins := os.Getenv("RESET_TOKEN_DURATION_MIN")
	if resetTokenDurationMins == "" {
		resetTokenDurationMins = "10"
	}
	mins, err := strconv.Atoi(resetTokenDurationMins)
	if err != nil {
		return fmt.Errorf("service.ForgotPassword: %w", err)
	}
	duration := time.Duration(mins) * time.Minute
	expiry := time.Now().Add(duration).Format(time.RFC3339)

	token, hashedToken, err := generatePasswordResetToken()
	if err != nil {
		return fmt.Errorf("service.ForgotPassword: %w", err)
	}

	resetURL := fmt.Sprintf("https://localhost:3000/execs/resetpassword/reset/%s", token)
	message := fmt.Sprintf("Forgot your password? Reset your password using the following link:\n"+
		"%s\n"+
		"If you didn't request a password reset, please ignore this email.\n"+
		"This link is only valid for %s minutes.",
		resetURL,
		resetTokenDurationMins,
	)

	err = sendEmail("schooladmin@school.com", email, "Password reset link", message)
	if err != nil {
		return fmt.Errorf("service.ForgotPassword: %w", err)
	}

	exec.PasswordResetToken = sql.NullString{String: hashedToken, Valid: true}
	exec.PasswordTokenExpires = sql.NullString{String: expiry, Valid: true}
	err = s.repo.UpdatePasswordResetToken(exec)
	if err != nil {
		return fmt.Errorf("service.ForgotPassword: %w", err)
	}
	return nil
}

func (s *ExecsService) ResetPassword(req *models.ResetPasswordRequest) error {
	if req.NewPassword == "" || req.ConfirmPassword == "" {
		return apperrors.NewError(apperrors.ErrValidation, errors.New("passwords are required"))
	}
	if req.NewPassword != req.ConfirmPassword {
		return apperrors.NewError(apperrors.ErrValidation, errors.New("passwords should match"))
	}
	encodedPassword, err := encodePassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("service.ResetPassword: %w", err)
	}

	bytes, err := hex.DecodeString(req.Token)
	if err != nil {
		return fmt.Errorf("service.ResetPassword: %w", err)
	}
	hashedTokenBytes := sha256.Sum256(bytes)
	hashedToken := hex.EncodeToString(hashedTokenBytes[:])

	exec, err := s.repo.GetByPasswordResetToken(hashedToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.NewError(apperrors.ErrNotFound, errors.New("password reset link is invalid or has already been used"))
		}
		return fmt.Errorf("service.ResetPassword: %w", err)
	}

	passwordTokenExpires, err := time.Parse(time.RFC3339, exec.PasswordTokenExpires.String)
	if err != nil {
		return fmt.Errorf("service.ResetPassword: %w", err)
	}
	if time.Now().After(passwordTokenExpires) {
		err = s.repo.ResetPassword(exec)
		if err != nil {
			return fmt.Errorf("service.ResetPassword: %w", err)
		}
		return apperrors.NewError(apperrors.ErrResourceGone, errors.New("password reset link has expired"))
	}

	currentTime := time.Now().Format(time.RFC3339)
	exec.PasswordChangedAt = sql.NullString{String: currentTime, Valid: true}
	exec.Password = encodedPassword

	err = s.repo.ResetPassword(exec)
	if err != nil {
		return fmt.Errorf("service.ResetPassword: %w", err)
	}
	return nil
}
