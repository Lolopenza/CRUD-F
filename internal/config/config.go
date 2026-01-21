package config

import (
	"fmt"
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

func Load() (*Config, error) {
	_ = godotenv.Load()

	db := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if err := validateDB(db); err != nil {
		return nil, err
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

	server := ServerConfig{
		Port: os.Getenv("SERVER_PORT"),
	}

	if server.Port == "" {
		return nil, fmt.Errorf("SERVER_PORT is required")
	}

	return &Config{
		Server: server,
		DB:     db,
	}, nil
}

func validateDB(db DBConfig) error {
	if db.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if db.Port == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if db.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if db.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if db.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if db.SSLMode == "" {
		return fmt.Errorf("DB_SSLMODE is required")
	}
	return nil
}
