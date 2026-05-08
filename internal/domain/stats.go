package domain

type AdminStats struct {
	TotalApplications int            `json:"total_applications"`
	TotalUsers        int            `json:"total_users"`
	TotalSessions     int            `json:"total_sessions"`
	Languages         []LanguageStat `json:"languages"`
}

type LanguageStat struct {
	Language string `json:"language"`
	Count    int    `json:"count"`
}
