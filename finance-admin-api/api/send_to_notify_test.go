package api

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func Test_parseNotifyApiKey(t *testing.T) {
	tests := []struct {
		name             string
		key              string
		expectedIss      string
		expectedJwtToken string
	}{
		{
			name:             "Empty API key",
			key:              "",
			expectedIss:      "",
			expectedJwtToken: "",
		},
		{
			name:             "API key with too many dashes",
			key:              "oh-no-1234abcd-1234-abcd-5678-123456abcdef-hehe0101-asdf-1234-hehe-12345678abcd",
			expectedIss:      "",
			expectedJwtToken: "",
		},
		{
			name:             "Normal shaped API key",
			key:              "hehe-1234abcd-1234-abcd-5678-123456abcdef-hehe0101-asdf-1234-hehe-12345678abcd",
			expectedIss:      "1234abcd-1234-abcd-5678-123456abcdef",
			expectedJwtToken: "hehe0101-asdf-1234-hehe-12345678abcd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iss, jwtToken := parseNotifyApiKey(tt.key)
			assert.Equal(t, tt.expectedIss, iss)
			assert.Equal(t, tt.expectedJwtToken, jwtToken)
		})
	}
}

func TestServer_SendEmailToNotify(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		expectedErr error
	}{
		{
			name:        "Status created",
			status:      http.StatusCreated,
			expectedErr: nil,
		},
		{
			name:   "Status unauthorized",
			status: http.StatusUnauthorized,
			expectedErr: StatusError{
				http.StatusUnauthorized,
				"//https:%2F%2Fapi.notifications.service.gov.uk",
				http.MethodPost,
			},
		},
		{
			name:   "Status OK",
			status: http.StatusOK,
			expectedErr: StatusError{
				http.StatusOK,
				"//https:%2F%2Fapi.notifications.service.gov.uk",
				http.MethodPost,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHttpClient := MockHttpClient{}
			server := Server{http: &mockHttpClient}
			ctx := context.Background()

			GetDoFunc = func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: tt.status,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
					Request: &http.Request{
						Method: http.MethodPost,
						URL:    &url.URL{Host: notifyUrl},
					},
				}, nil
			}

			err := server.SendEmailToNotify(ctx, "test@email.com", map[int]string{1: "test"}, "testReport")
			assert.Equal(t, tt.expectedErr, err)
		})
	}

}
