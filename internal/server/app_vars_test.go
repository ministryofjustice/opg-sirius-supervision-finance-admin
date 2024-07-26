package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewAppVars(t *testing.T) {
	r, _ := http.NewRequest("GET", "/path", nil)
	r.SetPathValue("clientId", "1")
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "abc123"})

	envVars := EnvironmentVars{}
	vars := NewAppVars(r, envVars)

	assert.Equal(t, AppVars{
		Path:            "/path",
		XSRFToken:       "abc123",
		EnvironmentVars: envVars,
		Tabs: []Tab{
			{
				Id:    "downloads",
				Title: "Downloads",
			},
			{
				Id:    "uploads",
				Title: "Uploads",
			},
			{
				Id:    "annual-invoicing-letters",
				Title: "Annual Invoicing Letters",
			},
		},
	}, vars)
}
