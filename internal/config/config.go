package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	DSN      string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using system env")
	}

	db := DBConfig{
		Host:     mustEnv("DB_HOST"),
		Port:     mustEnv("DB_PORT"),
		User:     mustEnv("DB_USER"),
		Password: mustEnv("DB_PASSWORD"),
		Name:     mustEnv("DB_NAME"),
		SSLMode:  mustEnv("DB_SSLMODE"),
	}

	db.DSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.SSLMode,
	)

	return &Config{
		Server: ServerConfig{
			Port: mustEnv("SERVER_PORT"),
		},
		DB: db,
	}
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required env variable: %s", key)
	}
	return value
}
