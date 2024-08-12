package server

import (
	"errors"
	"github.com/opg-sirius-finance-admin/internal/api"
	"github.com/opg-sirius-finance-admin/internal/model"
	"github.com/opg-sirius-finance-admin/internal/util/util"
	"net/http"
)

type DownloadHandler struct {
	router
}

func (h *DownloadHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)
	params := r.URL.Query()

	var (
		reportType         = params.Get("reportType")
		reportJournalType  = params.Get("reportJournalType")
		reportScheduleType = params.Get("reportScheduleType")
		reportAccountType  = params.Get("reportAccountType")
		reportDebtType     = params.Get("reportDebtType")
		dateField          = params.Get("dateField")
		dateFromField      = params.Get("dateFromField")
		dateToField        = params.Get("dateToField")
		emailField         = params.Get("emailField")
	)

	err := h.Client().Download(ctx, reportType, reportJournalType, reportScheduleType, reportAccountType, reportDebtType, dateField, dateFromField, dateToField, emailField)

	if err != nil {
		var (
			valErr model.ValidationError
			stErr  api.StatusError
		)
		if errors.As(err, &valErr) {
			data := AppVars{Errors: util.RenameErrors(valErr.Errors)}
			w.WriteHeader(http.StatusUnprocessableEntity)
			err = h.execute(w, r, data)
		} else if errors.As(err, &stErr) {
			data := AppVars{Error: stErr.Error()}
			w.WriteHeader(stErr.Code)
			err = h.execute(w, r, data)
		}
	}

	return err
}
