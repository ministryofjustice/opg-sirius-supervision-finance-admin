package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"github.com/stretchr/testify/assert"
)

func TestAnnualBillingLetters(t *testing.T) {
	client := &Client{
		http: &MockClient{},
		jwt:  &mockJWTClient{},
	}

	dataFromApiCall := shared.AnnualBillingInformation{
		AnnualBillingYear:        "",
		DemandedExpectedCount:    0,
		DemandedIssuedCount:      0,
		DemandedSkippedCount:     0,
		DirectDebitExpectedCount: 0,
		DirectDebitIssuedCount:   0,
		DirectDebitSkippedCount:  0,
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(dataFromApiCall)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		uid        string
		statusCode int
		wantErr    bool
		errType    error
	}{
		{
			name:       "successful request",
			uid:        "valid-uid",
			statusCode: http.StatusOK,
			wantErr:    false,
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
					Body:       io.NopCloser(&body),
				}, nil
			}

			_, err := client.AnnualBillingLetters(testContext())

			if tt.wantErr {
				assert.Error(t, err)
				assert.IsType(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
