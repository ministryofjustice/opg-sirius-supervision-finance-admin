package server

import (
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
)

type GetUploadsVars struct {
	ReportsUploadTypes *[]shared.ReportUploadType
	AppVars
}

type GetUploadsHandler struct {
	router
}

func (h *GetUploadsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetUploadsVars{&shared.ReportUploadTypes, v}
	data.selectTab("uploads")
	return h.execute(w, r, data)
}
