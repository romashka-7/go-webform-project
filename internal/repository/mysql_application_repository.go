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
		INSERT INTO applications (
		name,
		phone,
		email,
		birth_date,
		gender,
		biography,
		agreement
		)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`

	result, err := r.db.Exec(
		query,
		application.Name,
		application.Phone,
		application.Email,
		application.BirthDate,
		application.Gender,
		application.Biography,
		application.Agreement,
	)
	if err != nil {
		return domain.Application{}, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return domain.Application{}, err
	}

	application.ID = int(id)
	for _, languageID := range application.Languages {
		_, err := r.db.Exec(`
		INSERT INTO application_languages (
			application_id,
			language_id
		)
		VALUES (?, ?)
	`,
			application.ID,
			languageID,
		)

		if err != nil {
			return domain.Application{}, err
		}
	}

	return application, nil
}

func (r *MySQLApplicationRepository) GetAll() ([]domain.Application, error) {

	query := `
		SELECT
			id,
			name,
			COALESCE(phone, ''),
			email,
			COALESCE(birth_date, ''),
			COALESCE(gender, ''),
			COALESCE(biography, ''),
			agreement,
			created_at
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
			&application.Phone,
			&application.Email,
			&application.BirthDate,
			&application.Gender,
			&application.Biography,
			&application.Agreement,
			&application.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		languages, err := r.getApplicationLanguages(application.ID)

		if err != nil {
			return nil, err
		}

		application.Languages = languages
		applications = append(applications, application)
	}

	return applications, nil
}

func (r *MySQLApplicationRepository) Update(id int, application domain.Application) (domain.Application, error) {
	query := `
		UPDATE applications
SET
	name = ?,
	phone = ?,
	email = ?,
	birth_date = ?,
	gender = ?,
	biography = ?,
	agreement = ?
WHERE id = ?
	`

	result, err := r.db.Exec(
		query,
		application.Name,
		application.Phone,
		application.Email,
		application.BirthDate,
		application.Gender,
		application.Biography,
		application.Agreement,
		id,
	)
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
	_, err = r.db.Exec(`
	DELETE FROM application_languages
	WHERE application_id = ?
	`, id)

	for _, languageID := range application.Languages {
		_, err := r.db.Exec(`
		INSERT INTO application_languages (
			application_id,
			language_id
		)
		VALUES (?, ?)
	`,
			id,
			languageID,
		)

		if err != nil {
			return domain.Application{}, err
		}
	}

	if err != nil {
		return domain.Application{}, err
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

func (r *MySQLApplicationRepository) GetAdminStats() (domain.AdminStats, error) {
	var stats domain.AdminStats

	err := r.db.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM applications),
			(SELECT COUNT(*) FROM users),
			(SELECT COUNT(*) FROM sessions)
	`).Scan(
		&stats.TotalApplications,
		&stats.TotalUsers,
		&stats.TotalSessions,
	)

	if err != nil {
		return domain.AdminStats{}, err
	}

	rows, err := r.db.Query(`
		SELECT l.name, COUNT(al.application_id) AS count
		FROM languages l
		LEFT JOIN application_languages al ON al.language_id = l.id
		GROUP BY l.id, l.name
		ORDER BY count DESC, l.name ASC
	`)

	if err != nil {
		return domain.AdminStats{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var stat domain.LanguageStat

		err := rows.Scan(&stat.Language, &stat.Count)
		if err != nil {
			return domain.AdminStats{}, err
		}

		stats.Languages = append(stats.Languages, stat)
	}

	if err := rows.Err(); err != nil {
		return domain.AdminStats{}, err
	}

	return stats, nil
}

func (r *MySQLApplicationRepository) getApplicationLanguages(applicationID int) ([]int, error) {

	rows, err := r.db.Query(`
		SELECT language_id
		FROM application_languages
		WHERE application_id = ?
	`, applicationID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	languages := []int{}

	for rows.Next() {

		var languageID int

		err := rows.Scan(&languageID)
		if err != nil {
			return nil, err
		}

		languages = append(languages, languageID)
	}

	return languages, nil
}

func (r *MySQLApplicationRepository) GetByID(id int) (domain.Application, error) {
	query := `
		SELECT
			id,
			name,
			COALESCE(phone, ''),
			email,
			COALESCE(birth_date, ''),
			COALESCE(gender, ''),
			COALESCE(biography, ''),
			agreement,
			created_at
		FROM applications
		WHERE id = ?
	`

	var application domain.Application

	err := r.db.QueryRow(query, id).Scan(
		&application.ID,
		&application.Name,
		&application.Phone,
		&application.Email,
		&application.BirthDate,
		&application.Gender,
		&application.Biography,
		&application.Agreement,
		&application.CreatedAt,
	)

	if err != nil {
		return domain.Application{}, err
	}

	languages, err := r.getApplicationLanguages(application.ID)
	if err != nil {
		return domain.Application{}, err
	}

	application.Languages = languages

	return application, nil
}
