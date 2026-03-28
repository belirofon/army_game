package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	RedisURL    string
	GraphQLPath string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/army_game?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		GraphQLPath: getEnv("GRAPHQL_PATH", "/graphql"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
