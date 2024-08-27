package server

import (
	"github.com/opg-sirius-finance-admin/internal/model"
	"net/http"
)

type GetDownloadsVars struct {
	ReportsTypes        []model.ReportsType
	ReportJournalTypes  []model.ReportJournalType
	ReportScheduleTypes []model.ReportScheduleType
	ReportAccountTypes  []model.ReportAccountType
	ReportDebtTypes     []model.ReportDebtType
	AppVars
}

type GetDownloadsHandler struct {
	router
}

func (h *GetDownloadsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetDownloadsVars{
		model.ReportsTypes,
		model.ReportJournalTypes,
		model.ReportScheduleTypes,
		model.ReportAccountTypes,
		model.ReportDebtTypes,
		v}
	data.selectTab("downloads")
	return h.execute(w, r, data)
}
