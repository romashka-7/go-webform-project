package service

import (
	"webform-go/internal/domain"
	"webform-go/internal/repository"
)

type ApplicationService struct {
	repo repository.ApplicationRepository
}

func NewApplicationService(repo repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{
		repo: repo,
	}
}

func (s *ApplicationService) Create(application domain.Application) (domain.Application, error) {
	return s.repo.Save(application)
}
