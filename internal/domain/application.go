package domain

//package domain - modules of data that is structures which project trade inside themself and with json

type Application struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Biography string `json:"biography"`
	Agreement bool   `json:"agreement"`
	Languages []int  `json:"languages"`
	CreatedAt string `json:"created_at"`
}
