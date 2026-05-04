package repository

import "webform-go/internal/domain"

type MemoryApplicationRepository struct {
	applications []domain.Application
}

func NewMemoryApplicationRepository() *MemoryApplicationRepository {
	return &MemoryApplicationRepository{
		applications: []domain.Application{},
	}
}

func (r *MemoryApplicationRepository) Save(application domain.Application) (domain.Application, error) {
	r.applications = append(r.applications, application)

	return application, nil

}
