package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	ServerPort    string
	AdminUsername string
	AdminPassword string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment directly")
	}

	return Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "mirage"),
		ServerPort:    getEnv("PORT", getEnv("SERVER_PORT", "8080")),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
