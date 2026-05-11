package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
	AdminLogin    string
	AdminPassword string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	return Config{
		ServerPort: getEnv("PORT", getEnv("SERVER_PORT", "8080")),

		DBHost:     getEnv("MYSQLHOST", "localhost"),
		DBPort:     getEnv("MYSQLPORT", "3306"),
		DBUser:     getEnv("MYSQLUSER", "root"),
		DBPassword: getEnv("MYSQLPASSWORD", ""),
		DBName:     getEnv("MYSQLDATABASE", "go_webform"),

		AdminLogin:    getEnv("ADMIN_LOGIN", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}
	return value
}
