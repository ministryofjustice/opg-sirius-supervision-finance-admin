package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_CheckUserSession(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{
			name:   "valid session",
			status: http.StatusOK,
			want:   true,
		},
		{
			name:   "unauthorised",
			status: http.StatusUnauthorized,
			want:   false,
		},
		{
			name:   "server error",
			status: http.StatusInternalServerError,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			}))
			defer svr.Close()

			client, _ := NewClient(http.DefaultClient, svr.URL, "")

			got, _ := client.CheckUserSession(getContext(nil))
			assert.Equalf(t, tt.want, got, "CheckUserSession()")
		})
	}
}
