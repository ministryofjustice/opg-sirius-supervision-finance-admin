package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
)

type GetDownloadsVars struct {
	ReportJournalTypes  []model.ReportJournalType
	ReportScheduleTypes []model.ReportScheduleType
	ReportDebtTypes     []model.ReportDebtType
	ReportsTypes        []shared.ReportsType
	ReportAccountTypes  []shared.ReportAccountType
	AppVars
}

type DownloadsTabHandler struct {
	router
}

func (h *DownloadsTabHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetDownloadsVars{
		model.ReportJournalTypes,
		model.ReportScheduleTypes,
		model.ReportDebtTypes,
		shared.ReportsTypes,
		shared.ReportAccountTypes,
		v}
	data.selectTab("downloads")
	return h.execute(w, r, data)
}
