package services

import (
	"database/sql"
	"errors"
	"fmt"
	"restproject/internal/api/models"
	"restproject/internal/api/repositories"
	"restproject/internal/apperrors"
)

type TeachersService struct {
	repo *repositories.TeachersRepository
}

func NewTeachersService(repo *repositories.TeachersRepository) *TeachersService {
	return &TeachersService{repo: repo}
}

func (s *TeachersService) GetByID(id int) (*models.Teacher, error) {
	teacher, err := s.getByID(id)
	if err != nil {
		return nil, fmt.Errorf("service.GetByID: %w", err)
	}
	return teacher, nil
}

func (s *TeachersService) getByID(id int) (*models.Teacher, error) {
	teacher, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewError(apperrors.ErrNotFound, fmt.Errorf("teacher id %d not found", id))
		}
		return nil, err
	}
	return teacher, nil
}

func (s *TeachersService) GetAllByCriteria(criteria models.Criteria) ([]models.Teacher, error) {
	teachers, err := s.repo.GetAllByCriteria(criteria)
	if err != nil {
		return nil, fmt.Errorf("service.GetAllByCriteria: %w", err)
	}
	return teachers, nil
}

func (s *TeachersService) SaveAll(teachers []models.Teacher) ([]models.Teacher, error) {
	for _, teacher := range teachers {
		if err := checkBlankFields(teacher); err != nil {
			return nil, fmt.Errorf("service.SaveAll: %w", err)
		}
	}

	teachers, err := s.repo.SaveAll(teachers)
	if err != nil {
		return nil, fmt.Errorf("service.SaveAll: %w", err)
	}
	return teachers, nil
}

func (s *TeachersService) Replace(id int, updatedTeacher *models.Teacher) (*models.Teacher, error) {
	if err := checkBlankFields(*updatedTeacher); err != nil {
		return nil, fmt.Errorf("service.Replace: %w", err)
	}

	dbTeacher, err := s.getByID(id)
	if err != nil {
		return nil, fmt.Errorf("service.Replace: %w", err)
	}

	updatedTeacher.ID = dbTeacher.ID
	updatedTeacher, err = s.repo.Update(updatedTeacher)
	if err != nil {
		return nil, fmt.Errorf("service.Replace: %w", err)
	}
	return updatedTeacher, nil
}

func (s *TeachersService) Update(id int, update map[string]any) (*models.Teacher, error) {
	dbTeacher, err := s.getByID(id)
	if err != nil {
		return nil, fmt.Errorf("service.Update: %w", err)
	}

	if err = applyUpdates(dbTeacher, update); err != nil {
		return nil, fmt.Errorf("service.Update: %w", err)
	}

	updatedTeacher, err := s.repo.Update(dbTeacher)
	if err != nil {
		return nil, fmt.Errorf("service.Update: %w", err)
	}
	return updatedTeacher, nil
}

func (s *TeachersService) UpdateAll(updates []map[string]any) ([]models.Teacher, error) {
	updatedTeachers := make([]models.Teacher, 0, len(updates))

	for _, update := range updates {
		id, err := extractID(update)
		if err != nil {
			return nil, fmt.Errorf("service.UpdateAll: %w", err)
		}

		dbTeacher, err := s.getByID(id)
		if err != nil {
			return nil, fmt.Errorf("service.UpdateAll: %w", err)
		}

		if err = applyUpdates(dbTeacher, update); err != nil {
			return nil, fmt.Errorf("service.UpdateAll: %w", err)
		}

		updatedTeachers = append(updatedTeachers, *dbTeacher)
	}

	updatedTeachers, err := s.repo.UpdateAll(updatedTeachers)
	if err != nil {
		return nil, fmt.Errorf("service.UpdateAll: %w", err)
	}
	return updatedTeachers, nil
}

func (s *TeachersService) Delete(id int) error {
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

func (s *TeachersService) DeleteAll(ids []int) ([]int, error) {
	deletedIds, err := s.repo.DeleteAll(ids)
	if err != nil {
		return nil, fmt.Errorf("service.DeleteAll: %w", err)
	}
	return deletedIds, nil
}
