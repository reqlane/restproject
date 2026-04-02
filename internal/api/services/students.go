package services

import (
	"database/sql"
	"errors"
	"fmt"
	"restproject/internal/api/models"
	"restproject/internal/api/repositories"
	"restproject/internal/apperrors"
	"strings"
)

type StudentsService struct {
	repo *repositories.StudentsRepository
}

func NewStudentsService(repo *repositories.StudentsRepository) *StudentsService {
	return &StudentsService{repo: repo}
}

func (s *StudentsService) GetByID(id int) (*models.Student, error) {
	student, err := s.getByID(id)
	if err != nil {
		return nil, fmt.Errorf("service.GetByID: %w", err)
	}
	return student, nil
}

func (s *StudentsService) getByID(id int) (*models.Student, error) {
	student, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewError(apperrors.ErrNotFound, fmt.Errorf("student id %d not found", id))
		}
		return nil, err
	}
	return student, nil
}

func (s *StudentsService) GetAllByCriteria(criteria *models.Criteria, pg *models.Pagination) ([]models.Student, int, error) {
	students, totalCount, err := s.repo.GetAllByCriteria(criteria, pg)
	if err != nil {
		return nil, 0, fmt.Errorf("service.GetAllByCriteria: %w", err)
	}
	return students, totalCount, nil
}

func (s *StudentsService) SaveAll(students []models.Student) ([]models.Student, error) {
	for _, student := range students {
		if err := checkBlankFields(student); err != nil {
			return nil, fmt.Errorf("service.SaveAll: %w", err)
		}
	}

	students, err := s.repo.SaveAll(students)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1452") {
			return nil, apperrors.NewError(apperrors.ErrForeignKeyViolation, errors.New("class/class teacher does not exist"))
		}
		return nil, fmt.Errorf("service.SaveAll: %w", err)
	}
	return students, nil
}

func (s *StudentsService) Replace(id int, updatedStudent *models.Student) (*models.Student, error) {
	if err := checkBlankFields(*updatedStudent); err != nil {
		return nil, fmt.Errorf("service.Replace: %w", err)
	}

	dbStudent, err := s.getByID(id)
	if err != nil {
		return nil, fmt.Errorf("service.Replace: %w", err)
	}

	updatedStudent.ID = dbStudent.ID
	updatedStudent, err = s.repo.Update(updatedStudent)
	if err != nil {
		return nil, fmt.Errorf("service.Replace: %w", err)
	}
	return updatedStudent, nil
}

func (s *StudentsService) Update(id int, update map[string]any) (*models.Student, error) {
	dbStudent, err := s.getByID(id)
	if err != nil {
		return nil, fmt.Errorf("service.Update: %w", err)
	}

	if err = applyUpdates(dbStudent, update); err != nil {
		return nil, fmt.Errorf("service.Update: %w", err)
	}

	updatedStudent, err := s.repo.Update(dbStudent)
	if err != nil {
		return nil, fmt.Errorf("service.Update: %w", err)
	}
	return updatedStudent, nil
}

func (s *StudentsService) UpdateAll(updates []map[string]any) ([]models.Student, error) {
	updatedStudents := make([]models.Student, 0, len(updates))

	for _, update := range updates {
		id, err := extractID(update)
		if err != nil {
			return nil, fmt.Errorf("service.UpdateAll: %w", err)
		}

		dbStudent, err := s.getByID(id)
		if err != nil {
			return nil, fmt.Errorf("service.UpdateAll: %w", err)
		}

		if err = applyUpdates(dbStudent, update); err != nil {
			return nil, fmt.Errorf("service.UpdateAll: %w", err)
		}

		updatedStudents = append(updatedStudents, *dbStudent)
	}

	updatedStudents, err := s.repo.UpdateAll(updatedStudents)
	if err != nil {
		return nil, fmt.Errorf("service.UpdateAll: %w", err)
	}
	return updatedStudents, nil
}

func (s *StudentsService) Delete(id int) error {
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

func (s *StudentsService) DeleteAll(ids []int) ([]int, error) {
	deletedIds, err := s.repo.DeleteAll(ids)
	if err != nil {
		return nil, fmt.Errorf("service.DeleteAll: %w", err)
	}
	return deletedIds, nil
}
