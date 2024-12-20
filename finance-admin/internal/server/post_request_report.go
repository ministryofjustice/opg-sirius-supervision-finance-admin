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
		reportType         = params.Get("reportType")
		reportJournalType  = params.Get("reportJournalType")
		reportScheduleType = params.Get("reportScheduleType")
		reportAccountType  = params.Get("reportAccountType")
		reportDebtType     = params.Get("reportDebtType")
		dateOfTransaction  = params.Get("dateOfTransaction")
		dateFrom           = params.Get("dateFrom")
		dateTo             = params.Get("dateTo")
		email              = params.Get("email")
	)

	parsedReportAccountType := shared.ParseReportAccountType(reportAccountType)

	data := shared.NewReportRequest(reportType, reportJournalType, reportScheduleType, reportAccountType, reportDebtType, dateOfTransaction, dateTo, dateFrom, email)
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

	w.Header().Add("HX-Redirect", fmt.Sprintf("%s/downloads?success=request_report&reportAccountType=%s", v.EnvironmentVars.Prefix, parsedReportAccountType.Translation()))

	return err
}
