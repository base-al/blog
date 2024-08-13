// config/config.go
package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Database settings
	DBDriver   string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBPath     string
	DBURL      string

	// Server settings
	ServerPort string

	// Application settings
	Theme     string
	JWTSecret string
	AdminUser string
	AdminPass string
	Debug     bool
}

func New() *Config {
	return &Config{
		// Database settings
		DBDriver:   getEnv("DB_DRIVER", "sqlite"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "blog"),
		DBPath:     getEnv("DB_PATH", "blog.db"),
		DBURL:      getEnv("DB_URL", ""),

		// Server settings
		ServerPort: getEnv("PORT", "8080"),

		// Application settings
		Theme:     getEnv("THEME", "default"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		AdminUser: getEnv("ADMIN_USER", "admin"),
		AdminPass: getEnv("ADMIN_PASS", "password"),
		Debug:     getEnvAsBool("DEBUG", false),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	strValue := getEnv(key, "")
	if value, err := strconv.ParseBool(strValue); err == nil {
		return value
	}
	return fallback
}
