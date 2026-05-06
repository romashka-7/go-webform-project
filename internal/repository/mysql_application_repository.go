package repository

import (
	"database/sql"

	"webform-go/internal/domain"
)

type MySQLApplicationRepository struct {
	db *sql.DB
}

func NewMySQLApplicationRepository(db *sql.DB) *MySQLApplicationRepository {
	return &MySQLApplicationRepository{
		db: db,
	}
}

func (r *MySQLApplicationRepository) Save(application domain.Application) (domain.Application, error) {

	query := `
		INSERT INTO applications (name, email)
		VALUES (?, ?)
	`

	result, err := r.db.Exec(
		query,
		application.Name,
		application.Email,
	)

	if err != nil {
		return domain.Application{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return domain.Application{}, err
	}

	application.ID = int(id)

	return application, nil
}

func (r *MySQLApplicationRepository) GetAll() ([]domain.Application, error) {

	query := `
		SELECT id, name, email
		FROM applications
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var applications []domain.Application

	for rows.Next() {

		var application domain.Application

		err := rows.Scan(
			&application.ID,
			&application.Name,
			&application.Email,
		)

		if err != nil {
			return nil, err
		}

		applications = append(applications, application)
	}

	return applications, nil
}
