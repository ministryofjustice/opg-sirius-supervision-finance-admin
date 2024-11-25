package server

import (
	"github.com/opg-sirius-finance-admin/finance-admin/internal/components"
	"net/http"
)

type AnnualInvoicingLettersTabHandler struct {
	router
}

func (h *AnnualInvoicingLettersTabHandler) render(v components.AppVars, w http.ResponseWriter, r *http.Request) error {
	return h.execute(w, r, components.AnnualInvoicingLetters())
}
