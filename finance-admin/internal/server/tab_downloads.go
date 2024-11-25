package server

import (
	"github.com/opg-sirius-finance-admin/finance-admin/internal/components"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/model"
	"net/http"
)

type DownloadsTabHandler struct {
	router
}

func (h *DownloadsTabHandler) render(v components.AppVars, w http.ResponseWriter, r *http.Request) error {
	data := components.DownloadsTabVars{
		ReportsTypes:        model.ReportsTypes,
		ReportJournalTypes:  model.ReportJournalTypes,
		ReportScheduleTypes: model.ReportScheduleTypes,
		ReportAccountTypes:  model.ReportAccountTypes,
		ReportDebtTypes:     model.ReportDebtTypes,
		AppVars:             v}
	return h.execute(w, r, components.DownloadsTab(data))
}
