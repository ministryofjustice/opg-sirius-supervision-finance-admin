package server

import (
	"github.com/opg-sirius-finance-admin/internal/components"
	"net/http"
)

type GetDownloadsHandler struct {
	router
}

func (h *GetDownloadsHandler) render(v components.AppVars, w http.ResponseWriter, r *http.Request) error {
	return h.execute(w, r, components.Downloads())
}
