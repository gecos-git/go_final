package config

import (
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	Port   string
	DBfile string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Port:   getEnv("TODO_PORT", ":7540"),
		DBfile: getEnv("TODO_DBFILE", "scheduler.db"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
