package repository

import "webform-go/internal/domain"

type ApplicationRepository interface {
	Save(application domain.Application) (domain.Application, error)
}
