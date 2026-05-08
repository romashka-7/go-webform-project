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

	applications := []domain.Application{}

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

func (r *MySQLApplicationRepository) Update(id int, application domain.Application) (domain.Application, error) {
	query := `
		UPDATE applications
		SET name = ?, email = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, application.Name, application.Email, id)
	if err != nil {
		return domain.Application{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.Application{}, err
	}

	if rowsAffected == 0 {
		return domain.Application{}, sql.ErrNoRows
	}

	application.ID = id

	return application, nil
}

func (r *MySQLApplicationRepository) Delete(id int) error {
	query := `
		DELETE FROM applications
		WHERE id = ?
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *MySQLApplicationRepository) CreateUser(applicationID int, login string, passwordHash string) error {
	query := `
		INSERT INTO users (application_id, login, password_hash)
		VALUES (?, ?, ?)
	`

	_, err := r.db.Exec(query, applicationID, login, passwordHash)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLApplicationRepository) GetUserByLogin(login string) (domain.User, error) {
	query := `
		SELECT id, application_id, login, password_hash
		FROM users
		WHERE login = ?
	`

	var user domain.User

	err := r.db.QueryRow(query, login).Scan(
		&user.ID,
		&user.ApplicationID,
		&user.Login,
		&user.PasswordHash,
	)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *MySQLApplicationRepository) CreateSession(userID int, sessionID string) error {
	query := `
		INSERT INTO sessions (user_id, session_id)
		VALUES (?, ?)
	`

	_, err := r.db.Exec(query, userID, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLApplicationRepository) GetUserBySessionID(sessionID string) (domain.User, error) {
	query := `
		SELECT users.id, users.application_id, users.login, users.password_hash
		FROM users
		JOIN sessions ON sessions.user_id = users.id
		WHERE sessions.session_id = ?
	`

	var user domain.User

	err := r.db.QueryRow(query, sessionID).Scan(
		&user.ID,
		&user.ApplicationID,
		&user.Login,
		&user.PasswordHash,
	)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *MySQLApplicationRepository) DeleteSession(sessionID string) error {
	query := `
		DELETE FROM sessions
		WHERE session_id = ?
	`

	_, err := r.db.Exec(query, sessionID)
	if err != nil {
		return err
	}

	return nil
}
