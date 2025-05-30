package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/auth"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
)

type GetUploadsVars struct {
	UploadTypes []shared.ReportUploadType
	AppVars
}

type UploadsTabHandler struct {
	router
}

func (h *UploadsTabHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	var uploadTypes []shared.ReportUploadType

	user := r.Context().(auth.Context).User
	if user.IsCorporateFinance() {
		uploadTypes = append(uploadTypes, shared.PaymentUploadTypes...)
	}
	if user.IsFinanceReporting() {
		uploadTypes = append(uploadTypes, shared.ReportUploadTypes...)
	}

	data := GetUploadsVars{uploadTypes, v}
	data.selectTab("uploads")
	return h.execute(w, r, data)
}
