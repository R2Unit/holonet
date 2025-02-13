package config

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func PostgresFromEnv() Postgres {
	return Postgres{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnvAsInt("POSTGRES_PORT", 5432),
		User:     getEnv("POSTGRES_USER", "talos"),
		Password: getEnv("POSTGRES_PASSWORD", "saxofoon"),
		DBName:   getEnv("POSTGRES_DBNAME", "talos"),
	}
}
