package server

import (
	"net/http"
)

type GetAnnualInvoicingLettersVars struct {
	AppVars
}

type AnnualInvoicingLettersTabHandler struct {
	router
}

func (h *AnnualInvoicingLettersTabHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetAnnualInvoicingLettersVars{v}
	data.selectTab("annual-invoicing-letters")
	return h.execute(w, r, data)
}
