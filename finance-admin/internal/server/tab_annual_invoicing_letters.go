package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
)

type GetAnnualInvoicingLettersVars struct {
	AppVars
	shared.AnnualBillingInformation
}

type AnnualInvoicingLettersTabHandler struct {
	router
}

func (h *AnnualInvoicingLettersTabHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	annualBillingInfo, err := h.Client().AnnualBillingLetters(r.Context())
	if err != nil {
		log.Printf("Error calling download API: %v", err)
		http.Error(w, "Failed to call api", http.StatusInternalServerError)
	}
	fmt.Print("resp")
	fmt.Print(annualBillingInfo)
	data := GetAnnualInvoicingLettersVars{v, annualBillingInfo}
	data.selectTab("annual-invoicing-letters")
	return h.execute(w, r, data)
}
