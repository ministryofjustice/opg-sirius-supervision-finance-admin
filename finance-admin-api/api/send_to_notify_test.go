package api

import (
	"bytes"
	"context"
	"github.com/opg-sirius-finance-admin/shared"
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

func TestServer_formatFailedLines(t *testing.T) {
	tests := []struct {
		name        string
		failedLines map[int]string
		want        []string
	}{
		{
			name:        "Empty",
			failedLines: map[int]string{},
			want:        []string(nil),
		},
		{
			name: "Unsorted lines",
			failedLines: map[int]string{
				5: "DATE_PARSE_ERROR",
				3: "CLIENT_NOT_FOUND",
				8: "DUPLICATE_PAYMENT",
				1: "DUPLICATE_PAYMENT",
			},
			want: []string{
				"Line 1: Duplicate payment line",
				"Line 3: Could not find a client with this court reference",
				"Line 5: Unable to parse date",
				"Line 8: Duplicate payment line",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formattedLines := formatFailedLines(tt.failedLines)
			assert.Equal(t, tt.want, formattedLines)
		})
	}
}

func TestServer_createNotifyPayload(t *testing.T) {
	tests := []struct {
		name   string
		detail shared.FinanceAdminUploadProcessedEvent
		want   NotifyPayload
	}{
		{
			name: "Success",
			detail: shared.FinanceAdminUploadProcessedEvent{
				EmailAddress: "test@email.com",
				ReportType:   shared.ReportTypeUploadPaymentsMOTOCard.Key(),
			},
			want: NotifyPayload{
				EmailAddress: "test@email.com",
				TemplateId:   processingSuccessTemplateId,
				Personalisation: struct {
					ReportType string `json:"report_type"`
				}{shared.ReportTypeUploadPaymentsMOTOCard.Translation()},
			},
		},
		{
			name: "Failed",
			detail: shared.FinanceAdminUploadProcessedEvent{
				EmailAddress: "test@email.com",
				FailedLines: map[int]string{
					1: "DATE_PARSE_ERROR",
				},
				ReportType: shared.ReportTypeUploadPaymentsMOTOCard.Key(),
			},
			want: NotifyPayload{
				EmailAddress: "test@email.com",
				TemplateId:   processingFailedTemplateId,
				Personalisation: struct {
					FailedLines []string `json:"failed_lines"`
					ReportType  string   `json:"report_type"`
				}{[]string{"Line 1: Unable to parse date"}, shared.ReportTypeUploadPaymentsMOTOCard.Translation()},
			},
		},
		{
			name: "Errored",
			detail: shared.FinanceAdminUploadProcessedEvent{
				EmailAddress: "test@email.com",
				Error:        "Couldn't open report",
				ReportType:   shared.ReportTypeUploadPaymentsMOTOCard.Key(),
			},
			want: NotifyPayload{
				EmailAddress: "test@email.com",
				TemplateId:   processingErrorTemplateId,
				Personalisation: struct {
					Error      string `json:"error"`
					ReportType string `json:"report_type"`
				}{"Couldn't open report", shared.ReportTypeUploadPaymentsMOTOCard.Translation()},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := createNotifyPayload(tt.detail)
			assert.Equal(t, tt.want, payload)
		})
	}
}

func Test_SendEmailToNotify(t *testing.T) {
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

			detail := shared.FinanceAdminUploadProcessedEvent{
				EmailAddress: "test@email.com",
				FailedLines:  map[int]string{1: "test"},
				ReportType:   shared.ReportTypeUploadPaymentsMOTOCard.Key(),
			}

			err := server.SendEmailToNotify(ctx, detail)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
