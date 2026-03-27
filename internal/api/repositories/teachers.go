package repositories

import (
	"database/sql"
	"fmt"
	"restproject/internal/api/models"
	"strings"
)

type TeachersRepository struct {
	db *sql.DB
}

func NewTeacherRepository(db *sql.DB) *TeachersRepository {
	return &TeachersRepository{db: db}
}

func (r *TeachersRepository) GetByID(id int) (*models.Teacher, error) {
	var teacher models.Teacher
	query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id=?`
	err := r.db.QueryRow(query, id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err != nil {
		return nil, fmt.Errorf("repo.GetByID: %w", err)
	}
	return &teacher, nil
}

func (r *TeachersRepository) GetAllByCriteria(criteria models.Criteria) ([]models.Teacher, error) {
	var query strings.Builder
	query.WriteString(`SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1`)

	var args []any
	for dbField, value := range criteria.Filters {
		query.WriteString(" AND " + dbField + " =?")
		args = append(args, value)
	}

	addSorting(&query, criteria.Sortings, models.TeacherFieldNames)

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

func (r *TeachersRepository) SaveAll(teachers []models.Teacher) ([]models.Teacher, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo.SaveAll: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(generateInsertQuery("teachers", models.Teacher{}))
	if err != nil {
		return nil, fmt.Errorf("repo.SaveAll: %w", err)
	}
	defer stmt.Close()

	for i, newTeacher := range teachers {
		values := getStructValues(newTeacher)
		result, err := stmt.Exec(values...)
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

func (r *TeachersRepository) Update(teacher *models.Teacher) (*models.Teacher, error) {
	query := `UPDATE teachers SET first_name=?, last_name=?, email=?, class=?, subject=? WHERE id=?`
	_, err := r.db.Exec(query, teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject, teacher.ID)
	if err != nil {
		return nil, fmt.Errorf("repo.Update: %w", err)
	}
	return teacher, nil
}

func (r *TeachersRepository) UpdateAll(teachers []models.Teacher) ([]models.Teacher, error) {
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
		_, err := updateStmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject, teacher.ID)
		if err != nil {
			return nil, fmt.Errorf("repo.UpdateAll: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("repo.UpdateAll: %w", err)
	}
	return teachers, nil
}

func (r *TeachersRepository) Delete(id int) error {
	query := `DELETE FROM teachers WHERE id=?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	return nil
}

func (r *TeachersRepository) DeleteAll(ids []int) ([]int, error) {
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

func (r *TeachersRepository) GetStudentsByTeacherID(id int) ([]models.Student, error) {
	query := `SELECT s.id, s.first_name, s.last_name, s.email, s.class FROM students s JOIN teachers t ON s.class = t.class WHERE t.id=?`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("repo.GetStudentsByTeacherID: %w", err)
	}
	defer rows.Close()

	students := []models.Student{}
	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return nil, fmt.Errorf("repo.GetStudentsByTeacherID: %w", err)
		}
		students = append(students, student)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("repo.GetStudentsByTeacherID: %w", err)
	}
	return students, nil
}

func (r *TeachersRepository) GetStudentsCountByTeacherID(id int) (int, error) {
	var studentsCount int
	query := `SELECT COUNT(*) FROM students s JOIN teachers t ON s.class = t.class WHERE t.id=?`
	err := r.db.QueryRow(query, id).Scan(&studentsCount)
	if err != nil {
		return 0, fmt.Errorf("repo.GetStudentsCountByTeacherID: %w", err)
	}
	return studentsCount, nil
}
