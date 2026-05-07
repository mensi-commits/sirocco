package config

import "os"

type Config struct {
	Port string
	Token string
	ControllerURL string
	DataDir string
}

func Load() Config {
	return Config{
		Port: os.Getenv("PORT"),
		Token: os.Getenv("TOKEN"),
		ControllerURL: os.Getenv("CONTROLLER_URL"),
		DataDir: "/var/lib/sirocco",
	}
}