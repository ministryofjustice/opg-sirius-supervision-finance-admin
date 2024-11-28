package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
)

type GetDownloadsVars struct {
	ReportsTypes        []shared.ReportsType
	ReportJournalTypes  []shared.ReportJournalType
	ReportScheduleTypes []model.ReportScheduleType
	ReportDebtTypes     []model.ReportDebtType
	ReportAccountTypes  []shared.ReportAccountType
	AppVars
}

type DownloadsTabHandler struct {
	router
}

func (h *DownloadsTabHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetDownloadsVars{
		shared.ReportsTypes,
		shared.ReportJournalTypes,
		model.ReportScheduleTypes,
		model.ReportDebtTypes,
		shared.ReportAccountTypes,
		v}
	data.selectTab("downloads")
	return h.execute(w, r, data)
}
