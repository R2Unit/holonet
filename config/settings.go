package config

type Settings struct {
	Mode  string
	Debug string
}

func EnvToSettings() Settings {
	return Settings{
		Mode:  getEnv("MODE", "dev"),
		Debug: getEnv("DEBUG", "false"),
	}
}
