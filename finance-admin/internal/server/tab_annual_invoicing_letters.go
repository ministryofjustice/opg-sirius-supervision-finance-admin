package server

import (
	"github.com/opg-sirius-finance-admin/finance-admin/internal/components"
	"net/http"
)

type GetAnnualInvoicingLettersHandler struct {
	router
}

func (h *GetAnnualInvoicingLettersHandler) render(v components.AppVars, w http.ResponseWriter, r *http.Request) error {
	return h.execute(w, r, components.AnnualInvoicingLetters())
}
