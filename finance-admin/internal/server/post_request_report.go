package server

import (
	"errors"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/api"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
)

type RequestReportHandler struct {
	router
}

func (h *RequestReportHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	params := r.Form

	var (
		reportType             = params.Get("reportType")
		journalType            = params.Get("journalType")
		scheduleType           = params.Get("scheduleType")
		accountsReceivableType = params.Get("accountsReceivableType")
		debtType               = params.Get("debtType")
		transactionDate        = params.Get("transactionDate")
		dateFrom               = params.Get("dateFrom")
		dateTo                 = params.Get("dateTo")
		email                  = params.Get("email")
	)

	data := shared.NewReportRequest(reportType, journalType, scheduleType, accountsReceivableType, debtType, transactionDate, dateTo, dateFrom, email)
	err := h.Client().RequestReport(ctx, data)

	if err != nil {
		var (
			valErr model.ValidationError
			stErr  api.StatusError
		)
		if errors.As(err, &valErr) {
			data := AppVars{ValidationErrors: RenameErrors(valErr.Errors)}
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = h.execute(w, r, data)
		} else if errors.As(err, &stErr) {
			data := AppVars{Error: stErr.Error()}
			w.WriteHeader(stErr.Code)
			err = h.execute(w, r, data)
		}
	}

	var successMessage string
	switch data.ReportType {
	case shared.ReportsTypeAccountsReceivable:
		successMessage = data.AccountsReceivableType.Translation()
	case shared.ReportsTypeJournal:
		successMessage = data.JournalType.Translation()
	case shared.ReportsTypeSchedule:
		successMessage = data.ScheduleType.Translation()
	case shared.ReportsTypeDebt:
		successMessage = data.DebtType.Translation()
	default:
		successMessage = "UNKNOWN"
	}
	w.Header().Add("HX-Redirect", fmt.Sprintf("%s/downloads?success=request_report&report_type=%s", v.EnvironmentVars.Prefix, successMessage))

	return err
}
