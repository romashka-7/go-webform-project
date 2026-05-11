package repository

import "webform-go/internal/domain"

//interfaces for realisetion of repository methods

type ApplicationRepository interface {
	Save(application domain.Application) (domain.Application, error)

	GetAll() ([]domain.Application, error)

	GetByID(id int) (domain.Application, error)

	Update(id int, application domain.Application) (domain.Application, error)

	Delete(id int) error

	CreateUser(applicationID int, login string, passwordHash string) error
	GetUserByLogin(login string) (domain.User, error)

	CreateSession(userID int, sessionID string) error
	GetUserBySessionID(sessionID string) (domain.User, error)
	DeleteSession(sessionID string) error

	GetAdminStats() (domain.AdminStats, error)
}
