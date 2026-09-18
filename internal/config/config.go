package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppName    string
	AppPort    string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() Config {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return Config{
		AppName:    getEnv("APP_NAME", "agnos-backend"),
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "agnos"),
		DBPassword: getEnv("DB_PASSWORD", "agnos123"),
		DBName:     getEnv("DB_NAME", "agnos"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
