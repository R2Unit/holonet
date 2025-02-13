package config

import (
	"os"
	"strconv"
)

type Config struct {
	Settings Settings
	Postgres Postgres
	S3       S3
}

// (lorenzok) Check voor de fallback urls
// Fallback URLs zijn ook de Dev access for local testing
// Dus bij de dev compose niks aanpassen :p
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return fallback
}

// (lorenzo) Inladen van Configuratie naar cache en om te gebruiken voor componenten
// TODO: Robusteren handling maken voor failed states of errors...
func LoadConfig() Config {
	// Initializing Settings Environments
	settingsConfig := EnvToSettings()
	// Initializing Postgres Environments
	postgresConfig := PostgresFromEnv()
	// Initializing S3 Environments
	s3Config := NewS3FromEnv()

	return Config{
		Settings: settingsConfig,
		Postgres: postgresConfig,
		S3:       s3Config,
	}
}
