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
			name: "finance admin upload failed event",
			event: shared.Event{
				Source:     "opg.supervision.finance",
				DetailType: "finance-admin-upload-failed",
				Detail: shared.FinanceAdminUploadFailedEvent{
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
