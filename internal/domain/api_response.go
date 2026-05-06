package domain

type APIResponse struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
