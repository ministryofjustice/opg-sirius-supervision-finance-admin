package server

import (
	"os"
)

type EnvironmentVars struct {
	WebDir          string
	SiriusURL       string
	SiriusPublicURL string
	BackendURL      string
	HubURL          string
	Prefix          string
}

func NewEnvironmentVars() EnvironmentVars {
	return EnvironmentVars{
		WebDir:          os.Getenv("WEB_DIR"),
		SiriusURL:       os.Getenv("SIRIUS_URL"),
		SiriusPublicURL: os.Getenv("SIRIUS_PUBLIC_URL"),
		Prefix:          os.Getenv("PREFIX"),
		BackendURL:      os.Getenv("BACKEND_URL"),
		HubURL:          os.Getenv("HUB_URL"),
	}
}
