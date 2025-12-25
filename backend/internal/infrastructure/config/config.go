package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseDSN string
	ServerPort  string
}

func NewConfig() (*Config, error) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_DSN environment variable is required")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseDSN: dsn,
		ServerPort:  port,
	}, nil
}
