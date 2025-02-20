package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/auth"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewAppVars(t *testing.T) {
	r, _ := http.NewRequest("GET", "/path", nil)
	r = r.WithContext(auth.Context{XSRFToken: "abc123"})
	envVars := Envs{}
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
