package config

type Settings struct {
	Mode string
}

func EnvToSettings() Settings {
	return Settings{
		Mode: getEnv("MODE", "dev"),
	}
}
