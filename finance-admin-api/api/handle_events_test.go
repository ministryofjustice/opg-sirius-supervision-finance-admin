package api

import (
	"bytes"
	"encoding/json"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-finance-admin/apierror"
	"github.com/opg-sirius-finance-admin/shared"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_handleEvents(t *testing.T) {
	var e apierror.BadRequest
	tests := []struct {
		name        string
		event       shared.Event
		expectedErr error
	}{
		{
			name: "finance admin upload processed event",
			event: shared.Event{
				Source:     "opg.supervision.finance",
				DetailType: "finance-admin-upload-processed",
				Detail: shared.FinanceAdminUploadProcessedEvent{
					EmailAddress: "test@email.com",
					FailedLines: map[int]string{
						1: "DUPLICATE_PAYMENT",
					}},
			},
			expectedErr: nil,
		},
		{
			name: "unknown event",
			event: shared.Event{
				Source:     "opg.supervision.sirius",
				DetailType: "test",
			},
			expectedErr: e,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockHttpClient := MockHttpClient{}
			server := Server{http: &mockHttpClient}

			GetDoFunc = func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}, nil
			}

			var body bytes.Buffer
			_ = json.NewEncoder(&body).Encode(test.event)
			r := httptest.NewRequest(http.MethodPost, "/events", &body)
			ctx := telemetry.ContextWithLogger(r.Context(), telemetry.NewLogger("test"))
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			err := server.handleEvents(w, r)
			if test.expectedErr != nil {
				assert.ErrorAs(t, err, &test.expectedErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestServer_createUploadNotifyPayload(t *testing.T) {
	tests := []struct {
		name   string
		detail shared.FinanceAdminUploadProcessedEvent
		want   NotifyPayload
	}{
		{
			name: "Success",
			detail: shared.FinanceAdminUploadProcessedEvent{
				EmailAddress: "test@email.com",
				UploadType:   shared.ReportTypeUploadPaymentsMOTOCard.Key(),
			},
			want: NotifyPayload{
				EmailAddress: "test@email.com",
				TemplateId:   processingSuccessTemplateId,
				Personalisation: struct {
					UploadType string `json:"upload_type"`
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
				UploadType: shared.ReportTypeUploadPaymentsMOTOCard.Key(),
			},
			want: NotifyPayload{
				EmailAddress: "test@email.com",
				TemplateId:   processingFailedTemplateId,
				Personalisation: struct {
					FailedLines []string `json:"failed_lines"`
					UploadType  string   `json:"upload_type"`
				}{[]string{"Line 1: Unable to parse date"}, shared.ReportTypeUploadPaymentsMOTOCard.Translation()},
			},
		},
		{
			name: "Errored",
			detail: shared.FinanceAdminUploadProcessedEvent{
				EmailAddress: "test@email.com",
				Error:        "Couldn't open report",
				UploadType:   shared.ReportTypeUploadPaymentsMOTOCard.Key(),
			},
			want: NotifyPayload{
				EmailAddress: "test@email.com",
				TemplateId:   processingErrorTemplateId,
				Personalisation: struct {
					Error      string `json:"error"`
					UploadType string `json:"upload_type"`
				}{"Couldn't open report", shared.ReportTypeUploadPaymentsMOTOCard.Translation()},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := createUploadNotifyPayload(tt.detail)
			assert.Equal(t, tt.want, payload)
		})
	}
}
