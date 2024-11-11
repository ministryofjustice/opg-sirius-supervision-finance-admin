package server

import (
	"encoding/base64"
	"net/http"
)

type GetDownloadVars struct {
	Uid      string
	Filename string
	AppVars
}

type GetDownloadHandler struct {
	router
}

func (h *GetDownloadHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	uid := r.URL.Query().Get("uid")

	data := GetDownloadVars{Uid: uid, Filename: decryptFilename(uid), AppVars: v}
	return h.execute(w, r, data)
}

func decryptFilename(uid string) string {
	filename, _ := base64.StdEncoding.DecodeString(uid)
	return string(filename)
}
