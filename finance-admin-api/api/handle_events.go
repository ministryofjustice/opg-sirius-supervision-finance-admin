package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/apierror"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
)

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var event shared.Event
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		return apierror.BadRequestError("event", "unable to parse event", err)
	}

	if event.Source == shared.EventSourceFinanceHub && event.DetailType == shared.DetailTypeFinanceAdminUploadProcessed {
		if detail, ok := event.Detail.(shared.FinanceAdminUploadProcessedEvent); ok {

			payload := createUploadNotifyPayload(detail)
			err := s.SendEmailToNotify(ctx, payload)
			if err != nil {
				return err
			}
		}
	} else {
		return apierror.BadRequestError("event", fmt.Sprintf("could not match event: %s %s", event.Source, event.DetailType), errors.New("no match"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}

func createUploadNotifyPayload(detail shared.FinanceAdminUploadProcessedEvent) NotifyPayload {
	var payload NotifyPayload

	uploadType := shared.ParseReportUploadType(detail.UploadType)
	if detail.Error != "" {
		payload = NotifyPayload{
			detail.EmailAddress,
			processingErrorTemplateId,
			struct {
				Error      string `json:"error"`
				UploadType string `json:"upload_type"`
			}{
				detail.Error,
				uploadType.Translation(),
			},
		}
	} else if len(detail.FailedLines) != 0 {
		payload = NotifyPayload{
			detail.EmailAddress,
			processingFailedTemplateId,
			struct {
				FailedLines []string `json:"failed_lines"`
				UploadType  string   `json:"upload_type"`
			}{
				formatFailedLines(detail.FailedLines),
				uploadType.Translation(),
			},
		}
	} else {
		payload = NotifyPayload{
			detail.EmailAddress,
			processingSuccessTemplateId,
			struct {
				UploadType string `json:"upload_type"`
			}{uploadType.Translation()},
		}
	}

	return payload
}
