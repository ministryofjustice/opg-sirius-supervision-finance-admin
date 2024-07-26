package server

import (
	"net/http"
)

type GetDownloadsVars struct {
	AppVars
}

type GetDownloadsHandler struct {
	router
}

func (h *GetDownloadsHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetDownloadsVars{v}
	data.selectTab("downloads")
	return h.execute(w, r, data)
}
