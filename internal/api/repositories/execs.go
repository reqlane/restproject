package repositories

import (
	"database/sql"
	"fmt"
	"restproject/internal/api/models"
	"strings"
)

type ExecsRepository struct {
	db *sql.DB
}

func NewExecsRepository(db *sql.DB) *ExecsRepository {
	return &ExecsRepository{db: db}
}

func (r *ExecsRepository) GetByID(id int) (*models.Exec, error) {
	var exec models.Exec
	query := `SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE id=?`
	err := r.db.QueryRow(query, id).Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role)
	if err != nil {
		return nil, fmt.Errorf("repo.GetByID: %w", err)
	}
	return &exec, nil
}

func (r *ExecsRepository) GetByUsername(username string) (*models.Exec, error) {
	var exec models.Exec
	query := `SELECT id, first_name, last_name, email, username, password, user_created_at, inactive_status, role FROM execs WHERE username=?`
	err := r.db.QueryRow(query, username).Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role)
	if err != nil {
		return nil, fmt.Errorf("repo.GetByID: %w", err)
	}
	return &exec, nil
}

func (r *ExecsRepository) GetAllByCriteria(criteria models.Criteria) ([]models.Exec, error) {
	var query strings.Builder
	query.WriteString(`SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE 1=1`)

	var args []any
	for dbField, value := range criteria.Filters {
		query.WriteString(" AND " + dbField + " =?")
		args = append(args, value)
	}

	addSorting(&query, criteria.Sortings, models.ExecFieldNames)

	rows, err := r.db.Query(query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("repo.GetAllByCriteria: %w", err)
	}
	defer rows.Close()

	execs := make([]models.Exec, 0)
	for rows.Next() {
		var exec models.Exec
		err = rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role)
		if err != nil {
			return nil, fmt.Errorf("repo.GetAllByCriteria: %w", err)
		}
		execs = append(execs, exec)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("repo.GetAllByCriteria: %w", err)
	}
	return execs, nil
}

func (r *ExecsRepository) SaveAll(execs []models.Exec) ([]models.Exec, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo.SaveAll: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(generateInsertQuery("execs", models.Exec{}))
	if err != nil {
		return nil, fmt.Errorf("repo.SaveAll: %w", err)
	}
	defer stmt.Close()

	for i, newExec := range execs {
		values := getStructValues(newExec)
		result, err := stmt.Exec(values...)
		if err != nil {
			return nil, fmt.Errorf("repo.SaveAll: %w", err)
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("repo.SaveAll: %w", err)
		}
		execs[i].ID = int(lastID)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("repo.SaveAll: %w", err)
	}
	return execs, nil
}

func (r *ExecsRepository) Update(exec *models.Exec) (*models.Exec, error) {
	query := `UPDATE execs SET first_name=?, last_name=?, email=?, username=? WHERE id=?`
	_, err := r.db.Exec(query, exec.FirstName, exec.LastName, exec.Email, exec.Username, exec.ID)
	if err != nil {
		return nil, fmt.Errorf("repo.Update: %w", err)
	}
	return exec, nil
}

func (r *ExecsRepository) UpdateAll(execs []models.Exec) ([]models.Exec, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo.UpdateAll: %w", err)
	}
	defer tx.Rollback()

	updateStmt, err := tx.Prepare(`UPDATE execs SET first_name=?, last_name=?, email=?, username=? WHERE id=?`)
	if err != nil {
		return nil, fmt.Errorf("repo.UpdateAll: %w", err)
	}
	defer updateStmt.Close()

	for _, exec := range execs {
		_, err := updateStmt.Exec(exec.FirstName, exec.LastName, exec.Email, exec.Username, exec.ID)
		if err != nil {
			return nil, fmt.Errorf("repo.UpdateAll: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("repo.UpdateAll: %w", err)
	}
	return execs, nil
}

func (r *ExecsRepository) Delete(id int) error {
	query := `DELETE FROM execs WHERE id=?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	return nil
}
