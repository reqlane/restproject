package services

import "restproject/internal/api/repositories"

type ExecsService struct {
	repo *repositories.ExecsRepository
}

func NewExecsService(repo *repositories.ExecsRepository) *ExecsService {
	return &ExecsService{repo: repo}
}
