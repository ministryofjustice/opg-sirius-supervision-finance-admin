package server

import (
	"net/http"
)

type GetDownloadVars struct {
	Uid string
	AppVars
}

type GetDownloadHandler struct {
	router
}

func (h *GetDownloadHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetDownloadVars{r.URL.Query().Get("uid"), v}
	return h.execute(w, r, data)
}
