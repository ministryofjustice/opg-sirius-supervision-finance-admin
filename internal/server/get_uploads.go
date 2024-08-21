package server

import (
	"github.com/opg-sirius-finance-admin/internal/model"
	"net/http"
)

type GetUploadsVars struct {
	ReportsUploadTypes *[]model.ReportUploadType
	AppVars
}

type GetUploadsHandler struct {
	router
}

func (h *GetUploadsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetUploadsVars{&model.ReportUploadTypes, v}
	data.selectTab("uploads")
	return h.execute(w, r, data)
}
