package pg

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	ENV      = ".env"
	DB       = "POSTGRES_DB"
	PASSWORD = "POSTGRES_PASSWORD"
	USER     = "POSTGRES_USER"
	PORT     = "POSTGRES_PORT"
	HOST     = "POSTGRES_HOST"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

func newConfig() Config {
	godotenv.Load(ENV)
	return Config{
		Host:     envOrDefault(os.Getenv(HOST), "localhost"),
		Port:     envOrDefault(os.Getenv(PORT), "5432"),
		User:     envOrDefault(os.Getenv(USER), "opensky-collector"),
		Password: envOrDefault(os.Getenv(PASSWORD), "collector-pass"),
		Dbname:   envOrDefault(os.Getenv(DB), "telemetry"),
	}
}

func envOrDefault(variable string, def string) string {
	if variable == "" {
		return def
	}
	return variable
}
