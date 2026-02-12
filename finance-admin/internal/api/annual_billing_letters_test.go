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
		AnnualBillingYear:        "2025",
		DemandedExpectedCount:    13,
		DemandedIssuedCount:      1,
		DemandedSkippedCount:     3,
		DirectDebitExpectedCount: 55,
		DirectDebitIssuedCount:   3,
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(dataFromApiCall)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		uid          string
		statusCode   int
		wantErr      bool
		wantResponse shared.AnnualBillingInformation
		errType      error
	}{
		{
			name:       "successful request",
			uid:        "valid-uid",
			statusCode: http.StatusOK,
			wantErr:    false,
			wantResponse: shared.AnnualBillingInformation{
				AnnualBillingYear:        "2025",
				DemandedExpectedCount:    13,
				DemandedIssuedCount:      1,
				DemandedSkippedCount:     3,
				DirectDebitExpectedCount: 55,
				DirectDebitIssuedCount:   3,
			},
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

			resp, err := client.AnnualBillingLetters(testContext())

			if tt.wantErr {
				assert.Error(t, err)
				assert.IsType(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResponse, resp)
			}
		})
	}
}
