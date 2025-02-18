package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
)

type GetDownloadsVars struct {
	ReportsTypes                  []shared.ReportsType
	ReportJournalTypes            []shared.JournalType
	ReportScheduleTypes           []shared.ScheduleType
	ReportDebtTypes               []shared.DebtType
	ReportAccountsReceivableTypes []shared.AccountsReceivableType
	AppVars
}

type DownloadsTabHandler struct {
	router
}

func (h *DownloadsTabHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetDownloadsVars{
		shared.ReportsTypes,
		shared.JournalTypes,
		shared.ScheduleTypes,
		shared.DebtTypes,
		shared.AccountsReceivableTypes,
		v}
	data.selectTab("downloads")
	return h.execute(w, r, data)
}
