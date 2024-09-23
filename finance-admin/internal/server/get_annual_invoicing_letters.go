package server

import (
	"net/http"
)

type GetAnnualInvoicingLettersVars struct {
	AppVars
}

type GetAnnualInvoicingLettersHandler struct {
	router
}

func (h *GetAnnualInvoicingLettersHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	data := GetAnnualInvoicingLettersVars{v}
	data.selectTab("annual-invoicing-letters")
	return h.execute(w, r, data)
}
