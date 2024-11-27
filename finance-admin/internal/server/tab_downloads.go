package server

import (
	"github.com/opg-sirius-finance-admin/finance-admin/internal/model"
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
)

type GetDownloadsVars struct {
	ReportsTypes        []model.ReportsType
	ReportJournalTypes  []model.ReportJournalType
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
		model.ReportsTypes,
		model.ReportJournalTypes,
		model.ReportScheduleTypes,
		model.ReportDebtTypes,
		shared.ReportAccountTypes,
		v}
	data.selectTab("downloads")
	return h.execute(w, r, data)
}
