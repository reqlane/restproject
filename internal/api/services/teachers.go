package services

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"restproject/internal/api/models"
	"restproject/internal/api/repositories"
	"restproject/internal/apperrors"
	"strconv"
	"strings"
)

type TeachersService struct {
	repo *repositories.TeacherRepository
}

func NewTeachersService(repo *repositories.TeacherRepository) *TeachersService {
	return &TeachersService{repo: repo}
}

func (s *TeachersService) GetByID(id int) (*models.Teacher, error) {
	teacher, err := s.getByID(id)
	if err != nil {
		return nil, fmt.Errorf("service.getByID: %w", err)
	}
	return teacher, nil
}

func (s *TeachersService) getByID(id int) (*models.Teacher, error) {
	teacher, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return teacher, nil
}

func (s *TeachersService) GetAllByCriteria(criteria models.TeacherCriteria) ([]models.Teacher, error) {
	teachers, err := s.repo.GetAllByCriteria(criteria)
	if err != nil {
		return nil, fmt.Errorf("service.GetAllByCriteria: %w", err)
	}
	return teachers, nil
}

func (s *TeachersService) SaveAll(teachers []models.Teacher) ([]models.Teacher, error) {
	// TODO teachers fields validation
	teachers, err := s.repo.SaveAll(teachers)
	if err != nil {
		return nil, fmt.Errorf("service.SaveAll: %w", err)
	}
	return teachers, nil
}

func (s *TeachersService) Replace(id int, updatedTeacher *models.Teacher) (*models.Teacher, error) {
	// TODO updatedTeacher fields validation
	dbTeacher, err := s.getByID(id)
	if err != nil {
		return nil, fmt.Errorf("service.Replace: %w", err)
	}

	updatedTeacher.ID = dbTeacher.ID
	updatedTeacher, err = s.repo.Update(updatedTeacher)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("service.Replace: %w", apperrors.ErrNotFound)
		}
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("service.Update: %w", apperrors.ErrNotFound)
		}
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("service.UpdateAll: %w", apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("service.UpdateAll: %w", err)
	}
	return updatedTeachers, nil
}

func extractID(update map[string]any) (int, error) {
	idRaw, exists := update["id"]
	if !exists {
		return 0, apperrors.ErrMissingID
	}

	switch v := idRaw.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, apperrors.ErrInvalidID
		}
		return id, nil
	default:
		return 0, apperrors.ErrInvalidID
	}
}

func applyUpdates(teacher *models.Teacher, update map[string]any) error {
	teacherVal := reflect.ValueOf(teacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range update {
		if k == "id" {
			continue
		}
		for i := 0; i < teacherVal.NumField(); i++ {
			typeField := teacherType.Field(i)
			valField := teacherVal.Field(i)
			jsonName := strings.Split(typeField.Tag.Get("json"), ",")[0]
			if jsonName == k {
				if valField.CanSet() {
					value := reflect.ValueOf(v)
					if value.Type().ConvertibleTo(typeField.Type) {
						valField.Set(value.Convert(typeField.Type))
					} else {
						return apperrors.ErrInvalidField
					}
				}
				break
			}
		}
	}
	return nil
}

func (s *TeachersService) Delete(id int) error {
	err := s.repo.Delete(id)
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
