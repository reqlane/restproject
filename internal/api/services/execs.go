package services

import (
	"database/sql"
	"errors"
	"fmt"
	"restproject/internal/api/models"
	"restproject/internal/api/repositories"
	"restproject/internal/apperrors"
)

type ExecsService struct {
	repo *repositories.ExecsRepository
}

func NewExecsService(repo *repositories.ExecsRepository) *ExecsService {
	return &ExecsService{repo: repo}
}

func (s *ExecsService) GetByID(id int) (*models.Exec, error) {
	exec, err := s.getByID(id)
	if err != nil {
		return nil, fmt.Errorf("service.GetByID: %w", err)
	}
	return exec, nil
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

func (s *ExecsService) GetAllByCriteria(criteria models.Criteria) ([]models.Exec, error) {
	execs, err := s.repo.GetAllByCriteria(criteria)
	if err != nil {
		return nil, fmt.Errorf("service.GetAllByCriteria: %w", err)
	}
	return execs, nil
}

func (s *ExecsService) SaveAll(execs []models.Exec) ([]models.Exec, error) {
	for _, exec := range execs {
		if err := checkBlankFields(exec); err != nil {
			return nil, fmt.Errorf("service.SaveAll: %w", err)
		}
	}

	execs, err := s.repo.SaveAll(execs)
	if err != nil {
		return nil, fmt.Errorf("service.SaveAll: %w", err)
	}
	return execs, nil
}

func (s *ExecsService) Update(id int, update map[string]any) (*models.Exec, error) {
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
	return updatedExec, nil
}

func (s *ExecsService) UpdateAll(updates []map[string]any) ([]models.Exec, error) {
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
	return updatedExecs, nil
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
