package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opg-sirius-finance-admin/apierror"
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
)

func formatFailedLines(failedLines map[int]string) []string {
	var errorMessage string
	var formattedLines []string

	for i, line := range failedLines {
		switch line {
		case "DATE_PARSE_ERROR":
			errorMessage = "Unable to parse date"
		case "DUPLICATE":
			errorMessage = "Duplicate line"
		}
		formattedLines = append(formattedLines, fmt.Sprintf("Line %d: %s", i, errorMessage))
	}

	return formattedLines
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var event shared.Event
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		return apierror.BadRequestError("event", "unable to parse event", err)
	}

	if event.Source == shared.EventSourceFinanceHub && event.DetailType == shared.DetailTypeFinanceAdminUploadFailed {
		if detail, ok := event.Detail.(shared.FinanceAdminUploadFailedEvent); ok {
			templateId := "a8f9ab79-1489-4639-9e6c-cad1f079ebcf"

			err := s.SendEmailToNotify(ctx, detail.EmailAddress, templateId, formatFailedLines(detail.FailedLines), shared.ReportTypeUploadPaymentsMOTOCard.String())
			if err != nil {
				return err
			}
			fmt.Println(detail)
		}
	} else {
		return apierror.BadRequestError("event", fmt.Sprintf("could not match event: %s %s", event.Source, event.DetailType), errors.New("no match"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
