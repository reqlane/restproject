package repositories

import "database/sql"

type ExecsRepository struct {
	db *sql.DB
}

func NewExecsRepository(db *sql.DB) *ExecsRepository {
	return &ExecsRepository{db: db}
}
