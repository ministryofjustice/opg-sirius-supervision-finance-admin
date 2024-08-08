package server

import (
	"errors"
	"github.com/opg-sirius-finance-admin/internal/api"
	"github.com/opg-sirius-finance-admin/internal/model"
	"github.com/opg-sirius-finance-admin/internal/util/util"
	"net/http"
)

type SubmitDownloadHandler struct {
	router
}

func (h *SubmitDownloadHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	var (
		reportType         = r.PostFormValue("reportType")
		reportJournalType  = r.PostFormValue("reportJournalType")
		reportScheduleType = r.PostFormValue("reportScheduleType")
		reportAccountType  = r.PostFormValue("reportAccountType")
		reportDebtType     = r.PostFormValue("reportDebtType")
		dateField          = r.PostFormValue("dateField")
		dateFromField      = r.PostFormValue("dateFromField")
		dateToField        = r.PostFormValue("dateToField")
		emailField         = r.PostFormValue("emailField")
	)

	err := h.Client().SubmitDownload(ctx, reportType, reportJournalType, reportScheduleType, reportAccountType, reportDebtType, dateField, dateFromField, dateToField, emailField)

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
