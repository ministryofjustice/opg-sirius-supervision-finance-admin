package server

import (
	"github.com/opg-sirius-finance-admin/finance-admin/internal/components"
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
)

type UploadsTabHandler struct {
	router
}

func (h *UploadsTabHandler) render(v components.AppVars, w http.ResponseWriter, r *http.Request) error {
	data := components.UploadTabVars{ReportsUploadTypes: shared.ReportUploadTypes, AppVars: v}
	return h.execute(w, r, components.UploadsTab(data))
}
