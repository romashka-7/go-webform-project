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
		ServerPort: getEnv("SERVER_PORT", "8080"),

		DBUser:        getEnv("DB_USER", getEnv("MYSQLUSER", "root")),
		DBPassword:    getEnv("DB_PASSWORD", getEnv("MYSQLPASSWORD", "")),
		DBHost:        getEnv("DB_HOST", getEnv("MYSQLHOST", "localhost")),
		DBPort:        getEnv("DB_PORT", getEnv("MYSQLPORT", "3306")),
		DBName:        getEnv("DB_NAME", getEnv("MYSQLDATABASE", "go_webform")),
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
