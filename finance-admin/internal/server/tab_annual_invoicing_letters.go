package server

import (
	"fmt"
	"log"
	"net/http"
)

type GetAnnualInvoicingLettersVars struct {
	AppVars
}

type AnnualInvoicingLettersTabHandler struct {
	router
}

func (h *AnnualInvoicingLettersTabHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	resp, err := h.Client().AnnualBillingLetters(r.Context())
	if err != nil {
		log.Printf("Error calling download API: %v", err)
		http.Error(w, "Failed to stream file", http.StatusInternalServerError)
	}
	fmt.Print("resp")
	fmt.Print(resp.Body)
	data := GetAnnualInvoicingLettersVars{v}
	data.selectTab("annual-invoicing-letters")
	return h.execute(w, r, data)
}
