package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/auth"
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
	var reportsTypes []shared.ReportsType

	user := r.Context().(auth.Context).User
	if user.IsCorporateFinance() {
		reportsTypes = shared.PaymentReportsTypes
	}
	if user.IsFinanceReporting() {
		reportsTypes = shared.ReportsTypes
	}

	data := GetDownloadsVars{
		reportsTypes,
		shared.JournalTypes,
		shared.ScheduleTypes,
		shared.DebtTypes,
		shared.AccountsReceivableTypes,
		v}
	data.selectTab("downloads")
	return h.execute(w, r, data)
}
