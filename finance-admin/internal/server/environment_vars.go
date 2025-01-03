package server

import (
	"os"
)

type EnvironmentVars struct {
	Port            string
	WebDir          string
	SiriusURL       string
	SiriusPublicURL string
	BackendURL      string
	HubURL          string
	Prefix          string
}

func NewEnvironmentVars() EnvironmentVars {
	return EnvironmentVars{
		Port:            getEnv("PORT", "1234"),
		WebDir:          getEnv("WEB_DIR", "web"),
		SiriusURL:       getEnv("SIRIUS_URL", "http://host.docker.internal:8080"),
		SiriusPublicURL: getEnv("SIRIUS_PUBLIC_URL", ""),
		Prefix:          getEnv("PREFIX", ""),
		BackendURL:      getEnv("BACKEND_URL", ""),
		HubURL:          getEnv("HUB_URL", ""),
	}
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
