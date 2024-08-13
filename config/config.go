package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	Theme       string
	JWTSecret   string
	AdminUser   string
	AdminPass   string
	Debug       bool
}

func New() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "blog.db"),
		ServerPort:  getEnv("PORT", "8080"),
		Theme:       getEnv("THEME", "default"),
		JWTSecret:   getEnv("JWT_SECRET", "admin"),
		AdminUser:   getEnv("ADMIN_USER", "admin"),
		AdminPass:   getEnv("ADMIN_PASS", "password"),
		Debug:       getEnvAsBool("DEBUG", false),
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
