package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
)

type GetDownloadsVars struct {
	ReportsTypes                  []shared.ReportsType
	ReportJournalTypes            []shared.ReportJournalType
	ReportScheduleTypes           []shared.ReportScheduleType
	ReportDebtTypes               []shared.ReportDebtType
	ReportAccountsReceivableTypes []shared.ReportAccountsReceivableType
	AppVars
}

type DownloadsTabHandler struct {
	router
}

func (h *DownloadsTabHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetDownloadsVars{
		shared.ReportsTypes,
		shared.ReportJournalTypes,
		shared.ReportScheduleTypes,
		shared.ReportDebtTypes,
		shared.ReportAccountsReceivableTypes,
		v}
	data.selectTab("downloads")
	return h.execute(w, r, data)
}
