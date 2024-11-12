package server

import (
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
)

type GetDownloadVars struct {
	Uid      string `json:"uid"`
	Filename string `json:"filename"`
	AppVars
}

type GetDownloadHandler struct {
	router
}

func (h *GetDownloadHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	uid := r.URL.Query().Get("uid")

	var downloadRequest shared.DownloadRequest
	err := downloadRequest.Decode(uid)
	if err != nil {
		return err
	}
	data := GetDownloadVars{Uid: uid, Filename: downloadRequest.Key, AppVars: v}
	return h.execute(w, r, data)
}
