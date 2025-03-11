package config

type Settings struct {
	Mode  string
	Debug string
	Token string
}

func EnvToSettings() Settings {
	return Settings{
		Mode:  getEnv("MODE", "dev"),
		Debug: getEnv("DEBUG", "false"),
		Token: getEnv("TOKEN", "saxofoon"),
	}
}
