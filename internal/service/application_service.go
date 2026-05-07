package service

import (
	"errors"
	"webform-go/internal/domain"
	"webform-go/internal/repository"
	"webform-go/internal/security"
)

type ApplicationService struct {
	repo repository.ApplicationRepository
}

func NewApplicationService(repo repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{
		repo: repo,
	}
}

func (s *ApplicationService) Create(application domain.Application) (domain.Application, string, string, error) {
	createdApplication, err := s.repo.Save(application)
	if err != nil {
		return domain.Application{}, "", "", err
	}

	login := security.GenerateLogin()
	password := security.GeneratePassword()
	passwordHash := security.HashPassword(password)

	err = s.repo.CreateUser(createdApplication.ID, login, passwordHash)
	if err != nil {
		return domain.Application{}, "", "", err
	}

	return createdApplication, login, password, nil
}

func (s *ApplicationService) GetAll() ([]domain.Application, error) {
	return s.repo.GetAll()
}

func (s *ApplicationService) Update(id int, application domain.Application) (domain.Application, error) {
	return s.repo.Update(id, application)
}

func (s *ApplicationService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ApplicationService) Login(login string, password string) (domain.User, error) {
	user, err := s.repo.GetUserByLogin(login)
	if err != nil {
		return domain.User{}, err
	}

	if !security.CheckPassword(password, user.PasswordHash) {
		return domain.User{}, errors.New("invalid password")
	}

	return user, nil
}
