package api

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/apierror"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckDownload(t *testing.T) {
	client := &Client{
		http: &MockClient{},
		jwt:  &mockJWTClient{},
	}

	tests := []struct {
		name       string
		uid        string
		statusCode int
		wantErr    bool
		errType    error
	}{
		{
			name:       "successful download",
			uid:        "valid-uid",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "download not found",
			uid:        "invalid-uid",
			statusCode: http.StatusNotFound,
			wantErr:    true,
			errType:    apierror.NotFound{},
		},
		{
			name:       "internal server error",
			uid:        "server-error-uid",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
			errType:    StatusError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDoFunc = func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: tt.statusCode,
					Request:    req,
				}, nil
			}

			err := client.CheckDownload(testContext(), tt.uid)

			if tt.wantErr {
				assert.Error(t, err)
				assert.IsType(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
