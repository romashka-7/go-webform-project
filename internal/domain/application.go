package domain

type Application struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Languages []int  `json:"languages"`
	CreatedAt string `json:"created_at"`
}
