package config

import (
	"log"
	"os"
)

type AppConfig struct {
	BackendURL string
}

func NewAppConfig() AppConfig {
	url := os.Args
	if len(url) == 1 {
		log.Fatal("missing backend url")
	}
	return AppConfig{
		BackendURL: url[1],
	}
}
