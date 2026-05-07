package domain

type User struct {
	ID            int
	ApplicationID int
	Login         string
	PasswordHash  string
}
