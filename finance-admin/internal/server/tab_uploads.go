package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
)

type GetUploadsVars struct {
	ReportsUploadTypes []shared.ReportUploadType
	AppVars
}

type UploadsTabHandler struct {
	router
}

func (h *UploadsTabHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetUploadsVars{shared.ReportUploadTypes, v}
	data.selectTab("uploads")
	return h.execute(w, r, data)
}
