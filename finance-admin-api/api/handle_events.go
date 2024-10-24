package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opg-sirius-finance-admin/apierror"
	"github.com/opg-sirius-finance-admin/shared"
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
			err := s.SendEmailToNotify(ctx, detail, shared.ReportTypeUploadPaymentsMOTOCard.Translation())
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
