package repositories

import (
	"database/sql"
	"fmt"
	"restproject/internal/api/models"
	"restproject/internal/apperrors"
	"strings"
)

type TeacherRepository struct {
	db *sql.DB
}

func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) GetByID(id int) (*models.Teacher, error) {
	var teacher models.Teacher
	query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id=?`
	err := r.db.QueryRow(query, id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err != nil {
		return nil, fmt.Errorf("repo.GetByID: %w", err)
	}
	return &teacher, nil
}

func (r *TeacherRepository) GetAllByCriteria(criteria models.TeacherCriteria) ([]models.Teacher, error) {
	var query strings.Builder
	query.WriteString(`SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1`)

	var args []any
	for dbField, value := range criteria.Filters {
		query.WriteString(" AND " + dbField + " =?")
		args = append(args, value)
	}

	addSorting(&query, criteria.Sortings)

	rows, err := r.db.Query(query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("repo.GetAllByCriteria: %w", err)
	}
	defer rows.Close()

	teachers := make([]models.Teacher, 0)
	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return nil, fmt.Errorf("repo.GetAllByCriteria: %w", err)
		}
		teachers = append(teachers, teacher)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("repo.GetAllByCriteria: %w", err)
	}
	return teachers, nil
}

func addSorting(query *strings.Builder, sortings []string) {
	addedSort := false
	for _, param := range sortings {
		parts := strings.Split(param, ":")
		if len(parts) != 2 {
			continue
		}
		dbField, order := parts[0], parts[1]
		if !isValidSortField(dbField) || !isValidSortOrder(order) {
			continue
		}
		if !addedSort {
			query.WriteString(" ORDER BY")
			addedSort = true
		} else {
			query.WriteString(",")
		}
		query.WriteString(" " + dbField + " " + order)
	}
}

func isValidSortOrder(order string) bool {
	orderLowerCase := strings.ToLower(order)
	return orderLowerCase == "asc" || orderLowerCase == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

func (r *TeacherRepository) SaveAll(teachers []models.Teacher) ([]models.Teacher, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo.SaveAll: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)`)
	if err != nil {
		return nil, fmt.Errorf("repo.SaveAll: %w", err)
	}
	defer stmt.Close()

	for i, newTeacher := range teachers {
		result, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			return nil, fmt.Errorf("repo.SaveAll: %w", err)
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("repo.SaveAll: %w", err)
		}
		teachers[i].ID = int(lastID)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("repo.SaveAll: %w", err)
	}
	return teachers, nil
}

func (r *TeacherRepository) Update(teacher *models.Teacher) (*models.Teacher, error) {
	query := `UPDATE teachers SET first_name=?, last_name=?, email=?, class=?, subject=? WHERE id=?`
	result, err := r.db.Exec(query, teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject, teacher.ID)
	if err != nil {
		return nil, fmt.Errorf("repo.Update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("repo.Update: %w", err)
	}
	if rowsAffected < 1 {
		_, err := r.GetByID(teacher.ID)
		if err != nil {
			return nil, fmt.Errorf("repo.Update: %w", err)
		}
	}
	return teacher, nil
}

func (r *TeacherRepository) UpdateAll(teachers []models.Teacher) ([]models.Teacher, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo.UpdateAll: %w", err)
	}
	defer tx.Rollback()

	updateStmt, err := tx.Prepare(`UPDATE teachers SET first_name=?, last_name=?, email=?, class=?, subject=? WHERE id=?`)
	if err != nil {
		return nil, fmt.Errorf("repo.UpdateAll: %w", err)
	}
	defer updateStmt.Close()

	for _, teacher := range teachers {
		result, err := updateStmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject, teacher.ID)
		if err != nil {
			return nil, fmt.Errorf("repo.UpdateAll: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("repo.UpdateAll: %w", err)
		}
		if rowsAffected < 1 {
			_, err := r.GetByID(teacher.ID)
			if err != nil {
				return nil, fmt.Errorf("repo.UpdateAll: %w", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("repo.UpdateAll: %w", err)
	}
	return teachers, nil
}

func (r *TeacherRepository) Delete(id int) error {
	query := `DELETE FROM teachers WHERE id=?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	if rowsAffected < 1 {
		return fmt.Errorf("repo.Delete: %w", apperrors.ErrNotFound)
	}
	return nil
}

func (r *TeacherRepository) DeleteAll(ids []int) ([]int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo.DeleteAll: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM teachers WHERE id=?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("repo.DeleteAll: %w", err)
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			return nil, fmt.Errorf("repo.DeleteAll: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("repo.DeleteAll: %w", err)
		}

		if rowsAffected < 1 {
			return nil, fmt.Errorf("repo.DeleteAll: %w", apperrors.ErrNotFound)
		}
		deletedIds = append(deletedIds, id)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("repo.DeleteAll: %w", err)
	}
	return deletedIds, nil
}
