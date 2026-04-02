package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"restproject/internal/api/models"
	"strings"
)

type StudentsRepository struct {
	db *sql.DB
}

func NewStudentsRepository(db *sql.DB) *StudentsRepository {
	return &StudentsRepository{db: db}
}

func (r *StudentsRepository) GetByID(id int) (*models.Student, error) {
	var student models.Student
	query := `SELECT id, first_name, last_name, email, class FROM students WHERE id=?`
	err := r.db.QueryRow(query, id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err != nil {
		return nil, fmt.Errorf("repo.GetByID: %w", err)
	}
	return &student, nil
}

func (r *StudentsRepository) GetAllByCriteria(criteria *models.Criteria, pg *models.Pagination) ([]models.Student, int, error) {
	var query strings.Builder
	query.WriteString(`SELECT id, first_name, last_name, email, class FROM students WHERE 1=1`)

	var args []any
	for dbField, value := range criteria.Filters {
		query.WriteString(" AND " + dbField + " =?")
		args = append(args, value)
	}

	addSorting(&query, criteria.Sortings, models.StudentFieldNames)

	offset := (pg.Page - 1) * pg.Limit
	query.WriteString(" LIMIT ? OFFSET ?")
	args = append(args, pg.Limit, offset)

	rows, err := r.db.Query(query.String(), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("repo.GetAllByCriteria: %w", err)
	}
	defer rows.Close()

	students := make([]models.Student, 0)
	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return nil, 0, fmt.Errorf("repo.GetAllByCriteria: %w", err)
		}
		students = append(students, student)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, fmt.Errorf("repo.GetAllByCriteria: %w", err)
	}

	var totalCount int
	err = r.db.QueryRow(`SELECT COUNT(*) FROM students`).Scan(&totalCount)
	if err != nil {
		log.Println("repo.GetAllByCriteria:", err)
		totalCount = 0
	}
	return students, totalCount, nil
}

func (r *StudentsRepository) SaveAll(students []models.Student) ([]models.Student, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo.SaveAll: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(generateInsertQuery("students", models.Student{}))
	if err != nil {
		return nil, fmt.Errorf("repo.SaveAll: %w", err)
	}
	defer stmt.Close()

	for i, newStudent := range students {
		values := getStructValues(newStudent)
		result, err := stmt.Exec(values...)
		if err != nil {
			return nil, fmt.Errorf("repo.SaveAll: %w", err)
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("repo.SaveAll: %w", err)
		}
		students[i].ID = int(lastID)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("repo.SaveAll: %w", err)
	}
	return students, nil
}

func (r *StudentsRepository) Update(student *models.Student) (*models.Student, error) {
	query := `UPDATE students SET first_name=?, last_name=?, email=?, class=? WHERE id=?`
	_, err := r.db.Exec(query, student.FirstName, student.LastName, student.Email, student.Class, student.ID)
	if err != nil {
		return nil, fmt.Errorf("repo.Update: %w", err)
	}
	return student, nil
}

func (r *StudentsRepository) UpdateAll(students []models.Student) ([]models.Student, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo.UpdateAll: %w", err)
	}
	defer tx.Rollback()

	updateStmt, err := tx.Prepare(`UPDATE students SET first_name=?, last_name=?, email=?, class=? WHERE id=?`)
	if err != nil {
		return nil, fmt.Errorf("repo.UpdateAll: %w", err)
	}
	defer updateStmt.Close()

	for _, student := range students {
		_, err := updateStmt.Exec(student.FirstName, student.LastName, student.Email, student.Class, student.ID)
		if err != nil {
			return nil, fmt.Errorf("repo.UpdateAll: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("repo.UpdateAll: %w", err)
	}
	return students, nil
}

func (r *StudentsRepository) Delete(id int) error {
	query := `DELETE FROM students WHERE id=?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	return nil
}

func (r *StudentsRepository) DeleteAll(ids []int) ([]int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo.DeleteAll: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM students WHERE id=?`
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
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("repo.DeleteAll: %w", err)
	}
	return deletedIds, nil
}
